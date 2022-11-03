package handlers

import (
	"log"
	"net/http"

	"github.com/blazingly-fast/social-network/data"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Users struct {
	l *log.Logger
}

type KeyUser struct{}

type GenericError struct {
	Message string `json:"message"`
}

func NewUsers(l *log.Logger) *Users {
	return &Users{l}
}

func (u *Users) Signup(w http.ResponseWriter, r *http.Request) {

	user := &data.User{}

	err := data.FromJSON(user, r.Body)
	if err != nil {
		u.l.Println("[ERROR] deserializing product", err)

		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		u.l.Println("[ERROR] validating product", validationErr)

		w.WriteHeader(http.StatusUnprocessableEntity)
		data.ToJSON(&GenericError{Message: validationErr.Error()}, w)
		return
	}

	password := u.HashPassword(user.Password)
	user.Password = password

	//create token and append to user
	token, refreshToken, err := u.GenerateAllToken(user.Name, user.Password)
	user.Token = token
	user.Refresh_token = refreshToken

	data.Signup(user)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

}
