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
			a.l.Println("no token provided")
			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "no token provided"})
			return
		}
		claims, err := ValidateToken(clientToken)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, &GenericError{Message: err.Error()})
			return
		}

		if claims.UserType != "ADMIN" {
			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "unauthorized request!!!"})
			return
		}

		// isAdmin, err := a.store.FindAccountByUid(claims.Uid)
		// if err != nil {
		// 	WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "unauthorized request!!!"})
		// 	return
		// }

		// if isAdmin != true {
		// 	WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "unauthorized request!!!"})
		// 	return
		// }

		next.ServeHTTP(w, r)
	})
}
