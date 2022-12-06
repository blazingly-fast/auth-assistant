package handlers

import (
	"net/http"
)

func (s *Server) Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")
		if clientToken == "" {
			s.l.Println("no token provided")
			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "no token provided"})
			return
		}
		claims, err := ValidateToken(clientToken)
		if err != nil {
			s.l.Println(err)
			WriteJSON(w, http.StatusInternalServerError, &GenericError{Message: err.Error()})
			return
		}
		r.Header.Set("user_type", claims.UserType)
		r.Header.Set("uuid", claims.Uuid)
		r.Header.Set("email", claims.Email)
		next.ServeHTTP(w, r)
	})
}

// func (a *AccountHandler) IsAdmin(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		clientToken := r.Header.Get("token")
// 		if clientToken == "" {
// 			a.l.Println("no token provided")
// 			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "no token provided"})
// 			return
// 		}

// 		claims, err := ValidateToken(clientToken)
// 		if err != nil {
// 			a.l.Println(err)
// 			WriteJSON(w, http.StatusInternalServerError, &GenericError{Message: err.Error()})
// 			return
// 		}

// 		admin, err := a.store.GetAccountByField("uuid", claims.Uuid)
// 		if err != nil {
// 			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: err.Error()})
// 			return
// 		}

// 		if admin.UserType != "ADMIN" {
// 			a.l.Println("unauthorized request!!!")
// 			WriteJSON(w, http.StatusForbidden, &GenericError{Message: "unauthorized request!!!"})
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
