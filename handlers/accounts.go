package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/blazingly-fast/auth-assistant/data"
	"github.com/blazingly-fast/auth-assistant/util"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// HandleGetAccountByID handles GET request for single account
func (s *Server) HandleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	uuid := mux.Vars(r)["uuid"]

	if err := util.MatchUserTypeToUUID(r, uuid); err != nil {
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

// HandleGetAccounts handles GET requests and returns all current accounts
func (s *Server) HandleGetAccounts(w http.ResponseWriter, r *http.Request) error {

	if err := util.CheckUserType(r, "ADMIN"); err != nil {
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}

	pag := r.Context().Value(KeyHolder{}).(*Pagination)

	accountList, err := s.d.GetAccounts(pag.Limit, pag.CursorID)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accountList)
}

// HandleCreateAccount handles POST request to add new account
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

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	uuid := uuid.New().String()
	userType := "USER"
	avatar := "default.png"

	token, refreshToken, err := util.GenerateAllToken(
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
		avatar,
		uuid,
		token,
		refreshToken)

	err = s.d.CreateAccout(account)
	if err != nil {
		return err
	}

	res := data.NewAccountResponse(
		account.FirstName,
		account.LastName,
		account.Email,
		userType,
		avatar,
		uuid,
		token)

	return WriteJSON(w, http.StatusOK, &res)
}

// HandleUpdateAccount handles PUT/PATCH requests to update account
func (s *Server) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	req := &data.UpdateAccountRequest{}
	uuid := mux.Vars(r)["uuid"]

	if err := util.MatchUserTypeToUUID(r, uuid); err != nil {
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

	hashedPassword, err := util.HashPassword(req.Password)
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

// HandleDeleteAccount handles DELETE request to delete account
func (s *Server) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	uuid := mux.Vars(r)["uuid"]

	if err := util.CheckUserType(r, "ADMIN"); err != nil {
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

// HandleLogin handles POST login requests
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

	err = util.VerifyPassword(foundAccount.Password, req.Password)
	if err != nil {
		return err
	}

	token, refreshToken, _ := util.GenerateAllToken(
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
		foundAccount.Avatar,
		foundAccount.Uuid,
		token)

	return WriteJSON(w, http.StatusOK, &res)
}

// HandleAvatar handles POST request for the account avatar
func (s *Server) HandleAvatar(w http.ResponseWriter, r *http.Request) error {

	uuid := r.Header.Get("uuid")
	err := util.MatchUserTypeToUUID(r, uuid)
	if err != nil {
		s.l.Println("[ERROR]", err)
		return WriteJSON(w, http.StatusForbidden, &GenericError{Message: "Unauthorized to access this resource"})
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("upload_file")
	if err != nil {
		s.l.Println("[ERROR]", err)
		return WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "error retrieving file"})
	}
	defer file.Close()

	err = s.d.UpdateAvatar(handler.Filename, uuid)
	if err == data.ErrAccountNotFound {
		s.l.Println(err)
		return WriteJSON(w, http.StatusNotFound, &GenericError{Message: err.Error()})
	}
	if err != nil {
		return err
	}

	f, err := os.Create("./images/" + handler.Filename)
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(f, file)

	return WriteJSON(w, http.StatusOK, handler.Filename)
}
