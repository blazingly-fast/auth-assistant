package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/blazingly-fast/auth-assistant/data"
)

type Server struct {
	l *log.Logger
	v *data.Validation
	d data.Storer
}

func NewServer(l *log.Logger, v *data.Validation, d data.Storer) *Server {
	return &Server{
		l: l,
		v: v,
		d: d,
	}
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func (s *Server) MakeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := f(w, r); err != nil {
			s.l.Println(err)
			WriteJSON(w, http.StatusInternalServerError, &GenericError{Message: "Internal Server Error!"})
		}
	}
}

type GenericError struct {
	Message string `json:"message"`
}

type ValidationErrors struct {
	Messages []string `json:"messages"`
}

type Pagination struct {
	Limit    int
	CursorID int
}

type KeyHolder struct{}
