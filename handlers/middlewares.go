package handlers

import (
	"context"
	"net/http"
	"strconv"
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

const (
	// PageIDKey refers to the context key that stores the next page id
	PageIDKey CustomKey = "page_id"
)

type (
	CustomKey string
)

func (s *Server) Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageID := r.URL.Query().Get(string(PageIDKey))
		intPageID := 0
		var err error

		if pageID != "" {
			intPageID, err = strconv.Atoi(pageID)
			if err != nil {
				s.l.Println(err)
				WriteJSON(w, http.StatusBadRequest, &GenericError{Message: err.Error()})
				return
			}
		}

		ctx := context.WithValue(r.Context(), PageIDKey, intPageID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
