package main

import "net/http"

func (a *AccountHandler) Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")
		if clientToken == "" {
			a.l.Println("[ERROR] no token provided")
			WriteJSON(w, http.StatusBadRequest, "no token provided")
			return

		}

		claims, err := ValidateToken(clientToken)
		if err != "" {
			a.l.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			WriteJSON(w, http.StatusInternalServerError, "signature is invalid")
			return
		}

		r.Header.Set("name", claims.Email)
		next.ServeHTTP(w, r)
	})
}
