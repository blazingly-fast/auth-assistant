package main

import "time"

type CreateAccountRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	UserType  string `json:"user_type" validate:"required"`
	Uid       string `json:"uid" validate required`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Account struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name" validate:"required"`
	LastName     string    `json:"last_name" validate:"required"`
	Email        string    `json:"email" validate:"required"`
	Password     string    `json:"password" validate:"required"`
	UserType     string    `json:"user_type" validate:"required"`
	Uid          string    `json:"uid" validate:"required"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedOn    time.Time `json:"created_at"`
	UpdatedOn    time.Time `json:"updated_at"`
}

func NewAccount(firstName, lastName, email, password, userType, uid, token, refreshToken string) *Account {
	return &Account{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Password:     password,
		UserType:     userType,
		Uid:          uid,
		Token:        token,
		RefreshToken: refreshToken,
	}
}
