package main

import "net/http"

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
		a.l.Println("after")
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
		a.l.Println("is admin")
		next.ServeHTTP(w, r)
	})
}
