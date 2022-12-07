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

type Pagination struct {
	Page  int
	Limit int
}

func (s *Server) Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")
		pag := &Pagination{}
		var err error

		if limit != "" {
			pag.Limit, err = strconv.Atoi(limit)
			if err != nil {
				s.l.Println(err)
				WriteJSON(w, http.StatusBadRequest, &GenericError{Message: err.Error()})
				return
			}
		}

		if page != "" {
			pag.Page, err = strconv.Atoi(page)
			if err != nil {
				s.l.Println(err)
				WriteJSON(w, http.StatusBadRequest, &GenericError{Message: err.Error()})
				return
			}
		}

		ctx := context.WithValue(r.Context(), KeyHolder{}, pag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
