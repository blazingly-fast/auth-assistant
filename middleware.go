package main

import (
	"net/http"
)

func (a *AccountHandler) Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")
		if clientToken == "" {
			a.l.Println("no token provided")
			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "no token provided"})
			return
		}
		_, err := ValidateToken(clientToken)
		if err != nil {
			a.l.Println(err)
			WriteJSON(w, http.StatusInternalServerError, &GenericError{Message: err.Error()})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *AccountHandler) IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")

		if clientToken == "" {
			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "no token provided"})
			return
		}
		claims, err := ValidateToken(clientToken)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, &GenericError{Message: err.Error()})
			return
		}

		admin, err := a.store.FindAccountByUuid(claims.Uuid)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: err.Error()})
			return
		}

		if admin.UserType != "ADMIN" {
			a.l.Println("unauthorized request!!!")
			WriteJSON(w, http.StatusForbidden, &GenericError{Message: "unauthorized request!!!"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
