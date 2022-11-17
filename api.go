package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AccountHandler struct {
	l     *log.Logger
	store *PostgresStore
	v     *Validation
}

func NewAccountHandler(l *log.Logger, v *Validation, store *PostgresStore) *AccountHandler {
	return &AccountHandler{
		l:     l,
		store: store,
		v:     v,
	}
}

func (a *AccountHandler) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	acc, err := a.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (a *AccountHandler) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := a.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (a *AccountHandler) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	errs := a.v.Validate(req)
	if len(errs) != 0 {
		a.l.Println("[ERROR] validating request", errs)

		return WriteJSON(w, http.StatusUnprocessableEntity, &GenericErrors{Messages: errs.Errors()})
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}

	uuid := uuid.New().String()

	token, refreshToken, err := GenerateAllToken(
		req.FirstName,
		req.LastName,
		req.Email, uuid)

	account := NewAccount(
		req.FirstName,
		req.LastName,
		req.Email,
		hashedPassword,
		req.UserType,
		uuid,
		token, refreshToken)

	if err := a.store.CreateAccout(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, req)
}

func (a *AccountHandler) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	err = a.store.DeleteAccount(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (a *AccountHandler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	req := &LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	errs := a.v.Validate(req)
	if len(errs) != 0 {
		a.l.Println("[ERROR] validating request", errs)

		return WriteJSON(w, http.StatusUnprocessableEntity, &GenericErrors{Messages: errs.Errors()})
	}

	foundAccount, err := a.store.FindAccountByEmail(req.Email)
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
		foundAccount.Email, foundAccount.Uuid)

	err = a.store.UpdateAllTokens(token, refreshToken, foundAccount.ID)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, token)
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, err
	}

	return id, nil
}
