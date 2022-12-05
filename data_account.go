package main

import (
	"database/sql"
	"fmt"
	"time"
)

type Account struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name" validate:"required,min=2,max=50,alpha"`
	LastName     string    `json:"last_name" validate:"required,min=2,max=50,alpha"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password" validate:"required,min=8,max=50,containsany=1-9,containsany=Aa-Zz,alphanumunicode"`
	UserType     string    `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Uuid         string    `json:"uid" validate:"required,uuid"`
	Avatar       string    `json:"avatar"`
	Token        string    `json:"token" validate:"jwt"`
	RefreshToken string    `json:"refresh_token"`
	CreatedOn    time.Time `json:"created_at"`
	UpdatedOn    time.Time `json:"updated_at"`
}

func NewAccount(firstName, lastName, email, password, userType, uuid, avatar, token, refreshToken string) *Account {
	return &Account{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Password:     password,
		UserType:     userType,
		Uuid:         uuid,
		Avatar:       avatar,
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
	LastName  string    `json:"last_name" validate:"required,min=2,max=50,alpha"`
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

func (s *PostgresStore) CreateAccout(acc *Account) error {
	sql := `
	insert into account(first_name, last_name, email, password, user_type, uuid, avatar, token, refresh_token)
	values($1, $2, $3, $4, $5, $6, $7, $8, $9)
`
	_, err := s.db.Exec(
		sql, acc.FirstName,
		acc.LastName,
		acc.Email,
		acc.Password,
		acc.UserType,
		acc.Uuid,
		acc.Avatar,
		acc.Token,
		acc.RefreshToken)
	if err != nil {
		return err
	}

	return err
}

func (s *PostgresStore) UpdateAccount(acc *UpdateAccountRequest, uuid string) error {
	sql := `
	update account set 
	first_name=$1,
	last_name=$2,
	email=$3,
	password=$4,
	user_type=$5,
	updated_at=$6
	where uuid=$7
	`

	_, err := s.db.Exec(
		sql,
		acc.FirstName,
		acc.LastName,
		acc.Email,
		acc.Password,
		acc.UserType,
		acc.UpdatedOn,
		uuid)
	if err != nil {
		return err
	}

	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, err
}

func (s *PostgresStore) GetAccountByField(field string, value any) (*Account, error) {
	sql := fmt.Sprintf("select * from account where %s=$1", field)
	rows, err := s.db.Query(sql, value)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, ErrAccountNotFound
}

func (s *PostgresStore) DeleteAccount(uuid string) error {
	rows, err := s.db.Exec("delete from account where uuid = $1", uuid)
	count, _ := rows.RowsAffected()
	if count != 1 {
		return ErrAccountNotFound
	}

	return err
}

func (s *PostgresStore) UpdateAvatar(avatarURL, uuid string) error {
	rows, err := s.db.Exec("update account set avatar=$1 where uuid=$2", avatarURL, uuid)
	if err != nil {
		return err
	}

	count, _ := rows.RowsAffected()
	if count != 1 {
		return ErrAccountNotFound
	}

	return err
}

func (s *PostgresStore) UpdateAllTokens(token string, refreshToken string, id int) error {
	rows, err := s.db.Query("select * from account where ID=$1", id)
	if err != nil {
		return err
	}

	acc := &Account{}
	for rows.Next() {
		acc, err = scanIntoAccount(rows)
	}
	if err != nil {
		return err
	}

	acc.Token = token
	acc.RefreshToken = refreshToken

	sql := fmt.Sprintf("update account set token='%s', refresh_token='%s' where id=%d", acc.Token, acc.RefreshToken, id)

	_, err = s.db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	acc := &Account{}
	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Password,
		&acc.Email,
		&acc.UserType,
		&acc.Uuid,
		&acc.Token,
		&acc.Avatar,
		&acc.RefreshToken,
		&acc.CreatedOn,
		&acc.UpdatedOn,
	)
	return acc, err
}
