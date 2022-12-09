package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/blazingly-fast/social-network/util"
)

func (s *Server) Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")
		if clientToken == "" {
			s.l.Println("no token provided")
			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "no token provided"})
			return
		}
		claims, err := util.ValidateToken(clientToken)
		if err != nil {
			s.l.Println(err)
			WriteJSON(w, http.StatusInternalServerError, &GenericError{Message: "Internal Server Error!"})
			return
		}
		r.Header.Set("user_type", claims.UserType)
		r.Header.Set("email", claims.Email)
		r.Header.Set("uuid", claims.Uuid)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil && limitStr != "" {
			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "limit parameter is invalid"})
			return
		}
		if limit == 0 {
			limit = 10
		}

		cursorStr := r.URL.Query().Get("cursor")
		cursor, err := strconv.Atoi(cursorStr)
		if err != nil && cursorStr != "" {
			WriteJSON(w, http.StatusBadRequest, &GenericError{Message: "cursor parameter is invalid"})
			return
		}

		pag := &Pagination{
			Limit:    limit,
			CursorID: cursor,
		}

		ctx := context.WithValue(r.Context(), KeyHolder{}, pag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
