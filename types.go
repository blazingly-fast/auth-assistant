package main

import "time"

type CreateAccountRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=2,max=50"`
	UserType  string `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=2,max=50"`
}

type Account struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name" validate:"required"`
	LastName     string    `json:"last_name" validate:"required"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password" validate:"required,min=2,max=50"`
	UserType     string    `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Uuid         string    `json:"uid" validate:"required"`
	Token        string    `json:"token"`
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

type GenericError struct {
	Message string `json:"message"`
}

type GenericErrors struct {
	Messages []string `json:"messages"`
}

type KeyAccount struct{}
