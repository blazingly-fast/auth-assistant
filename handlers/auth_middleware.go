package handlers

import (
	"net/http"

	"github.com/blazingly-fast/social-network/data"
)

func (u *Users) Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")
		if clientToken == "" {
			u.l.Println("[ERROR] no token provided")

			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: "no token provided"}, w)
			return

		}

		claims, err := ValidateToken(clientToken)
		if err != "" {
			u.l.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&GenericError{Message: "signature is invalid"}, w)
			return
		}

		r.Header.Set("name", claims.Name)
		next.ServeHTTP(w, r)
	})
}
