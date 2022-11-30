package main

import "time"

type Account struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name" validate:"required,alpha"`
	LastName     string    `json:"last_name" validate:"required,alpha"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password" validate:"required,min=8,max=50,containsany=1-9,containsany=Aa-Zz,alphanumunicode"`
	UserType     string    `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Uuid         string    `json:"uid" validate:"required,uuid"`
	Token        string    `json:"token" validate:"jwt"`
	RefreshToken string    `json:"refresh_token"`
	CreatedOn    time.Time `json:"created_at"`
	UpdatedOn    time.Time `json:"updated_at"`
}

func NewAccount(firstName, lastName, email, password, userType, uuid, token, refreshToken string) *Account {
	return &Account{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Password:     password,
		UserType:     userType,
		Uuid:         uuid,
		Token:        token,
		RefreshToken: refreshToken,
	}
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50,alpha"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50,alpha"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=50,containsany=1-9,containsany=Aa-Zz,alphanumunicode"`
}

type UpdateAccountRequest struct {
	FirstName string    `json:"first_name" validate:"required,min=2,max=50,alpha"`
	LastName  string    `json:"last_name" validate:"required,min=2,max=50,aplpha"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=8,max=50,containsany=1-9,containsany=Aa-Zz,alphanumunicode"`
	UserType  string    `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	UpdatedOn time.Time `json:"updated_at" validate:"required"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=50,containsany=1-9,containsany=Aa-Zz,alphanumunicode"`
}

type AccountResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	UseryType string `json:"user_type"`
	Uuid      string `json:"uuid"`
	Token     string `json:"token"`
}

func NewAccountResponse(firstName, lastName, email, userType, uuid, token string) *AccountResponse {
	return &AccountResponse{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		UseryType: userType,
		Uuid:      uuid,
		Token:     token,
	}
}

type GenericError struct {
	Message string `json:"message"`
}

type ValidationErrors struct {
	Messages []string `json:"messages"`
}

type KeyAccount struct{}
