package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/blazingly-fast/social-network/data"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
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
		u.l.Println("[ERROR] deserializing user", err)

		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		u.l.Println("[ERROR] validating user", validationErr)

		w.WriteHeader(http.StatusUnprocessableEntity)
		data.ToJSON(&GenericError{Message: validationErr.Error()}, w)
		return
	}

	password := u.HashPassword(user.Password)
	user.Password = password

	//create token and append to user
	token, refreshToken, err := u.GenerateAllToken(user.Name)
	user.Token = token
	user.Refresh_token = refreshToken

	data.Signup(user)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	var user data.User

	err := data.FromJSON(&user, r.Body)
	if err != nil {
		u.l.Println("[ERROR] deserializing user", err)
		w.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	found_user := data.CheckEmail(&user)

	passwordValid, msg := u.VerifyPassword(user.Password, found_user.Password)
	if passwordValid != true {
		w.WriteHeader(http.StatusBadRequest)
		u.l.Println("[ERROR] invalid password", msg)
		data.ToJSON(&GenericError{Message: msg}, w)
		return
	}

	token, refreshToken, _ := u.GenerateAllToken(found_user.Name)

	data.UpdateAllTokens(token, refreshToken, found_user.ID)

	w.WriteHeader(http.StatusOK)
	data.ToJSON(&found_user, w)
	return
}

func (u *Users) VerifyPassword(pass string, hashedPass string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg

}

func (u *Users) GetUsers(w http.ResponseWriter, r *http.Request) {
	u.l.Println("true")
}
