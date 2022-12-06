package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/blazingly-fast/social-network/data"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) HandleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	uuid := mux.Vars(r)["uuid"]

	if err := MatchUserTypeToUUID(r, uuid); err != nil {
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}

	acc, err := s.d.GetAccountByField("uuid", uuid)
	if err == data.ErrAccountNotFound {
		return WriteJSON(w, http.StatusNotFound, &GenericError{Message: err.Error()})
	}
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (s *Server) HandleGetAccounts(w http.ResponseWriter, r *http.Request) error {

	if err := CheckUserType(r, "ADMIN"); err != nil {
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}
	// pageID := r.Context().Value(PageIDKey)

	accounts, err := s.d.GetAccounts(2)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *Server) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := &data.CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.l.Println(err)
		return err
	}

	errs := s.v.Validate(req)
	if len(errs) != 0 {
		s.l.Println("[ERROR] validating request", errs)
		return WriteJSON(w, http.StatusUnprocessableEntity, &ValidationErrors{Messages: errs.Errors()})
	}

	exists, _ := s.d.GetAccountByField("email", req.Email)
	if exists != nil {
		return WriteJSON(w, http.StatusBadRequest, &GenericError{Message: fmt.Sprintf("email %s already exists", req.Email)})
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}

	uuid := uuid.New().String()
	userType := "USER"
	avatar := "default.png"

	token, refreshToken, err := GenerateAllToken(
		req.FirstName,
		req.LastName,
		req.Email,
		userType,
		uuid)

	account := data.NewAccount(
		req.FirstName,
		req.LastName,
		req.Email,
		hashedPassword,
		userType,
		uuid,
		avatar,
		token,
		refreshToken)

	if err := s.d.CreateAccout(account); err != nil {
		return err
	}

	res := data.NewAccountResponse(
		account.FirstName,
		account.LastName,
		account.Email,
		userType,
		uuid,
		token)

	return WriteJSON(w, http.StatusOK, &res)
}

func (s *Server) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	req := &data.UpdateAccountRequest{}
	uuid := mux.Vars(r)["uuid"]

	if err := MatchUserTypeToUUID(r, uuid); err != nil {
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	req.UpdatedOn = time.Now().UTC()

	errs := s.v.Validate(req)
	if len(errs) != 0 {
		s.l.Println("[ERROR] validating request", errs)
		return WriteJSON(w, http.StatusUnprocessableEntity, &ValidationErrors{Messages: errs.Errors()})
	}

	foundAccWithUUID, err := s.d.GetAccountByField("uuid", uuid)
	if err == data.ErrAccountNotFound {
		return WriteJSON(w, http.StatusNotFound, &GenericError{Message: err.Error()})
	}

	foundAccWithEmail, err := s.d.GetAccountByField("email", req.Email)

	if foundAccWithEmail != nil && foundAccWithUUID.Email != req.Email {
		return WriteJSON(w, http.StatusUnprocessableEntity, &GenericError{Message: fmt.Sprintf("email %s already exists", req.Email)})
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}
	req.Password = hashedPassword

	err = s.d.UpdateAccount(req, uuid)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, fmt.Sprintf("account updated successfully"))
}

func (s *Server) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	uuid := mux.Vars(r)["uuid"]

	if err := CheckUserType(r, "ADMIN"); err != nil {
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}

	err := s.d.DeleteAccount(uuid)
	if err == data.ErrAccountNotFound {
		return WriteJSON(w, http.StatusNotFound, &GenericError{Message: err.Error()})
	}
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"deleted": uuid})
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	req := &data.LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	errs := s.v.Validate(req)
	if len(errs) != 0 {
		s.l.Println("[ERROR] validating request", errs)
		return WriteJSON(w, http.StatusUnprocessableEntity, &ValidationErrors{Messages: errs.Errors()})
	}

	foundAccount, err := s.d.GetAccountByField("email", req.Email)
	if err == data.ErrAccountNotFound {
		return WriteJSON(w, http.StatusNotFound, &GenericError{Message: err.Error()})
	}
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

	err = s.d.UpdateAllTokens(token, refreshToken, foundAccount.ID)
	if err != nil {
		return err
	}

	res := data.NewAccountResponse(
		foundAccount.FirstName,
		foundAccount.LastName,
		foundAccount.Email,
		foundAccount.UserType,
		foundAccount.Uuid,
		token)

	return WriteJSON(w, http.StatusOK, &res)
}

func (s *Server) HandleAvatar(w http.ResponseWriter, r *http.Request) error {

	uuid := r.Header.Get("uuid")
	err := MatchUserTypeToUUID(r, uuid)
	if err != nil {
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("upload_file")
	if err != nil {
		return err
	}
	defer file.Close()

	err = s.d.UpdateAvatar(handler.Filename, uuid)
	if err == data.ErrAccountNotFound {
		return WriteJSON(w, http.StatusNotFound, &GenericError{Message: err.Error()})
	}
	if err != nil {
		return err
	}

	f, err := os.OpenFile("./images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(f, file)

	return WriteJSON(w, http.StatusOK, handler.Filename)
}
