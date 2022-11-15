package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

func (h *AccountHandler) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	acc, err := h.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, acc)
}

func (h *AccountHandler) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	req := r.Context().Value(KeyAccount{}).(*CreateAccountRequest)
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return err
	}

	token, refreshToken, err := GenerateAllToken(req.FirstName, req.LastName, req.Email, req.UserType, req.Uid)

	account := NewAccount(
		req.FirstName,
		req.LastName,
		req.Email,
		hashedPassword,
		req.UserType,
		req.Uid,
		token, refreshToken)

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

	foundAccount, err := h.store.FindAccountByEmail(req)
	if err != nil {
		return err
	}

	if foundAccount.ID == 0 {
		return WriteJSON(w, http.StatusBadRequest, "email does not exist")
	}

	err = VerifyPassword(foundAccount.Password, req.Password)
	if err != nil {
		return err
	}

	token, refreshToken, _ := GenerateAllToken(
		foundAccount.FirstName,
		foundAccount.LastName,
		foundAccount.Email,
		foundAccount.UserType, foundAccount.Uid)

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

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, err
	}

	return id, nil
}
