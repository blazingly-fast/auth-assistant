package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

type AccountHandler struct {
	l     *log.Logger
	store *PostgresStore
}

func NewAccountHandler(l *log.Logger, store *PostgresStore) *AccountHandler {
	return &AccountHandler{
		l:     l,
		store: store,
	}
}

func (h *AccountHandler) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id := getID(r)
	acc, err := h.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (h *AccountHandler) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	if err := validate.Struct(req); err != nil {
		return err
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}

	token, refreshToken, err := GenerateAllToken(req.FirstName, req.LastName, req.Email)

	account := NewAccount(req.FirstName, req.LastName, req.Email, hashedPassword, token, refreshToken)

	if err := h.store.CreateAccout(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, req)
}

func (h *AccountHandler) handleLogin(w http.ResponseWriter, r *http.Request) error {
	req := &LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	if err := validate.Struct(req); err != nil {
		return err
	}

	foundAccount, err := h.store.CheckEmail(req)
	if err != nil {
		return err
	}

	err = VerifyPassword(foundAccount.Password, req.Password)
	if err != nil {
		return err
	}

	token, refreshToken, _ := GenerateAllToken(foundAccount.FirstName, foundAccount.LastName, foundAccount.Email)

	err = h.store.UpdateAllTokens(token, refreshToken, foundAccount.ID)
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

func getID(r *http.Request) int {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	return id
}
