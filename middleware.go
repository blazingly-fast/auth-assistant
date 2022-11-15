package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func (a *AccountHandler) Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")
		if clientToken == "" {
			a.l.Println("no token provided")
			WriteJSON(w, http.StatusBadRequest, "no token provided")
			return
		}

		claims, msg := ValidateToken(clientToken)
		if msg != "" {
			a.l.Println(msg)
			w.WriteHeader(http.StatusInternalServerError)
			WriteJSON(w, http.StatusInternalServerError, "signature is invalid")
			return
		}

		r.Header.Set("first_name", claims.FirstName)
		r.Header.Set("last_name", claims.LastName)
		r.Header.Set("email", claims.Email)
		r.Header.Set("user_type", claims.UserType)
		r.Header.Set("uid", claims.UserType)
		next.ServeHTTP(w, r)
	})
}

func (a *AccountHandler) IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")

		claims, _ := ValidateToken(clientToken)
		if claims.UserType != "ADMIN" {
			WriteJSON(w, http.StatusBadRequest, "Unauthorized request")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *AccountHandler) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &CreateAccountRequest{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteJSON(w, http.StatusBadRequest, err)
			return
		}

		errs := a.v.Validate(req)
		if len(errs) != 0 {
			a.l.Println("[ERROR] validating request", errs)

			WriteJSON(w, http.StatusUnprocessableEntity, &GenericErrors{Messages: errs.Errors()})
			return

		}

		ctx := context.WithValue(r.Context(), KeyAccount{}, req)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
