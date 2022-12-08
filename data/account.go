package data

import (
	"database/sql"
	"fmt"
	"time"
)

// Account defines the structure for an API account
type Account struct {
	ID           int       `json:"id"`
	FirstName    string    `json:"first_name" validate:"required,min=2,max=50,alpha"`
	LastName     string    `json:"last_name" validate:"required,min=2,max=50,alpha"`
	Email        string    `json:"email" validate:"required,email"`
	Password     string    `json:"password" validate:"required,min=8,max=50,containsany=1-9,containsany=Aa-Zz,alphanumunicode"`
	UserType     string    `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Avatar       string    `json:"avatar"`
	Uuid         string    `json:"uid" validate:"required,uuid"`
	Token        string    `json:"token" validate:"jwt"`
	RefreshToken string    `json:"refresh_token"`
	CreatedOn    time.Time `json:"created_at"`
	UpdatedOn    time.Time `json:"updated_at"`
}

func NewAccount(firstName, lastName, email, password, userType, avatar, uuid, token, refreshToken string) *Account {
	return &Account{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Password:     password,
		UserType:     userType,
		Avatar:       avatar,
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
	Avatar    string `json:"avatar"`
	Uuid      string `json:"uuid"`
	Token     string `json:"token"`
}

func NewAccountResponse(firstName, lastName, email, userType, avatar, uuid, token string) *AccountResponse {
	return &AccountResponse{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		UseryType: userType,
		Avatar:    avatar,
		Uuid:      uuid,
		Token:     token,
	}
}

type AccountList struct {
	Accounts []*Account `json:"accounts"`
	CursorID int        `json:"cursor_id,omitempty" exapmle:"10"`
}

var ErrAccountNotFound = fmt.Errorf("Account not found")

func (s *PostgresStore) CreateAccout(acc *Account) error {
	sql := `
	insert into account(first_name, last_name, email, password, user_type, avatar, uuid, token, refresh_token)
	values($1, $2, $3, $4, $5, $6, $7, $8, $9)
`
	_, err := s.db.Exec(
		sql, acc.FirstName,
		acc.LastName,
		acc.Email,
		acc.Password,
		acc.UserType,
		acc.Avatar,
		acc.Uuid,
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

func (s *PostgresStore) GetAccounts(limit, cursorID int) (*AccountList, error) {

	rows, err := s.db.Query("select * from account where id > $1 order by id limit $2", cursorID, limit)
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

	lastID := 0
	if len(accounts) > 0 {
		lastID = accounts[len(accounts)-1].ID
	}

	accList := &AccountList{
		Accounts: accounts,
		CursorID: lastID,
	}

	return accList, err
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
	rows, err := s.db.Query("select * from account where id=$1", id)
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

	_, err = s.db.Exec("update account set token=$1, refresh_token=$2 where id=$3", acc.Token, acc.RefreshToken, id)
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
		&acc.Avatar,
		&acc.Uuid,
		&acc.Token,
		&acc.RefreshToken,
		&acc.CreatedOn,
		&acc.UpdatedOn,
	)
	return acc, err
}
