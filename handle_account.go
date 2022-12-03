package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AccountHandler struct {
	l     *log.Logger
	v     *Validation
	store Storer
}

func NewAccountHandler(l *log.Logger, v *Validation, store Storer) *AccountHandler {
	return &AccountHandler{
		l:     l,
		v:     v,
		store: store,
	}
}

func (h *AccountHandler) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	uuid := mux.Vars(r)["uuid"]

	if err := MatchUserTypeToUUID(r, uuid); err != nil {
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}

	acc, err := h.store.GetAccountByField("uuid", uuid)
	if err == ErrAccountNotFound {
		return WriteJSON(w, http.StatusNotFound, &GenericError{Message: ErrAccountNotFound.Error()})
	}
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (h *AccountHandler) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {

	if err := CheckUserType(r, "ADMIN"); err != nil {
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}
	accounts, err := h.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (h *AccountHandler) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Println(err)
		return err
	}

	errs := h.v.Validate(req)
	if len(errs) != 0 {
		h.l.Println("[ERROR] validating request", errs)
		return WriteJSON(w, http.StatusUnprocessableEntity, &ValidationErrors{Messages: errs.Errors()})
	}

	exists, _ := h.store.GetAccountByField("email", req.Email)
	if exists != nil {
		return WriteJSON(w, http.StatusBadRequest, &GenericError{Message: fmt.Sprintf("email %s already exists", req.Email)})
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}

	uuid := uuid.New().String()
	userType := "USER"

	token, refreshToken, err := GenerateAllToken(
		req.FirstName,
		req.LastName,
		req.Email,
		userType,
		uuid)

	account := NewAccount(
		req.FirstName,
		req.LastName,
		req.Email,
		hashedPassword,
		userType,
		uuid,
		token,
		refreshToken)

	if err := h.store.CreateAccout(account); err != nil {
		return err
	}

	res := NewAccountResponse(
		account.FirstName,
		account.LastName,
		account.Email,
		userType,
		uuid,
		token)

	return WriteJSON(w, http.StatusOK, &res)
}

func (h *AccountHandler) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	req := &UpdateAccountRequest{}
	uuid := mux.Vars(r)["uuid"]

	if err := MatchUserTypeToUUID(r, uuid); err != nil {
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	req.UpdatedOn = time.Now().UTC()

	errs := h.v.Validate(req)
	if len(errs) != 0 {
		h.l.Println("[ERROR] validating request", errs)
		return WriteJSON(w, http.StatusUnprocessableEntity, &ValidationErrors{Messages: errs.Errors()})
	}

	foundAccWithUUID, err := h.store.GetAccountByField("uuid", uuid)
	if err == ErrAccountNotFound {
		return WriteJSON(w, http.StatusNotFound, &GenericError{Message: ErrAccountNotFound.Error()})
	}

	foundAccWithEmail, err := h.store.GetAccountByField("email", req.Email)

	if foundAccWithEmail != nil && foundAccWithUUID.Email != req.Email {
		return WriteJSON(w, http.StatusUnprocessableEntity, &GenericError{Message: fmt.Sprintf("email %s already exists", req.Email)})
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}
	req.Password = hashedPassword

	err = h.store.UpdateAccount(req, uuid)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, fmt.Sprintf("account updated successfully"))
}

func (h *AccountHandler) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	uuid := mux.Vars(r)["uuid"]

	if err := CheckUserType(r, "ADMIN"); err != nil {
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}

	err := h.store.DeleteAccount(uuid)
	if err == ErrAccountNotFound {
		return WriteJSON(w, http.StatusNotFound, &GenericError{Message: ErrAccountNotFound.Error()})
	}
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"deleted": uuid})
}

func (h *AccountHandler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	req := &LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	errs := h.v.Validate(req)
	if len(errs) != 0 {
		h.l.Println("[ERROR] validating request", errs)
		return WriteJSON(w, http.StatusUnprocessableEntity, &ValidationErrors{Messages: errs.Errors()})
	}

	foundAccount, err := h.store.GetAccountByField("email", req.Email)
	if err != nil {
		return err
	}

	err = VerifyPassword(foundAccount.Password, req.Password)
	if err != nil {
		return err
	}

	token, refreshToken, _ := GenerateAllToken(
		foundAccount.FirstName,
		foundAccount.LastName,
		foundAccount.Email,
		foundAccount.UserType,
		foundAccount.Uuid)

	err = h.store.UpdateAllTokens(token, refreshToken, foundAccount.ID)
	if err != nil {
		return err
	}

	res := NewAccountResponse(
		foundAccount.FirstName,
		foundAccount.LastName,
		foundAccount.Email,
		foundAccount.UserType,
		foundAccount.Uuid,
		token)

	return WriteJSON(w, http.StatusOK, &res)
}

func (h *AccountHandler) handleAvatar(w http.ResponseWriter, r *http.Request) error {

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("upload_file")
	if err != nil {
		h.l.Println("yo")
		return err
	}
	defer file.Close()

	f, err := os.OpenFile("./images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		h.l.Println("ye")
		h.l.Println(handler.Filename)
		return err
	}
	defer f.Close()
	io.Copy(f, file)

	return WriteJSON(w, http.StatusOK, handler.Filename)
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func CheckUserType(r *http.Request, role string) error {
	userType := r.Header.Get("user_type")

	if userType != role {
		return fmt.Errorf("Unauthorized to access this resource")
	}

	return nil
}

func MatchUserTypeToUUID(r *http.Request, claimsUUID string) error {
	userType := r.Header.Get("user_type")
	uuid := r.Header.Get("uuid")

	if userType != "ADMIN" && uuid != claimsUUID {
		return fmt.Errorf("Unauthorized to access this resource")
	}
	err := CheckUserType(r, userType)
	return err
}
