package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "miro"
	password = "bta"
	dbname   = "miro"
)

type Storage interface {
	// CreateAccout(*Account) error
	// DeleteAccount(int) error
	// UpdateAccount(*Account) error
	// GetAccounts() ([]*Account, error)
	// GetAccountByID(int) (*Account, error)
	// CheckEmail(*Account) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {

	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) createAccountTable() error {
	createSql := `
	  create table if not exists accounts(
	  id SERIAL PRIMARY KEY,
	  first_name text,
	  last_name text,
	  password text,
      email text,	
	  user_type text,
      uuid text,
	  token text,
	  refresh_token text,
	  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),	
	  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	  );
	  `
	_, err := s.db.Exec(createSql)
	return err
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) CreateAccout(acc *Account) error {
	sql := `
	insert into accounts(first_name, last_name, email, password, user_type, uuid, token, refresh_token)
	values($1, $2, $3, $4, $5, $6, $7, $8)
`
	_, err := s.db.Exec(
		sql, acc.FirstName,
		acc.LastName,
		acc.Email,
		acc.Password,
		acc.UserType,
		acc.Uuid,
		acc.Token, acc.RefreshToken)
	if err != nil {
		return err
	}

	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from accounts where id = $1", id)
	return err
}

func (s *PostgresStore) UpdateAccount(acc *UpdateAccountRequest, id int) error {
	sql := `
	update accounts set 
	first_name=$1,
	last_name=$2,
	email=$3,
	password=$4,
	user_type=$5,
	updated_at=$6
	where id=$7
	`

	_, err := s.db.Exec(
		sql,
		acc.FirstName,
		acc.LastName,
		acc.Email,
		acc.Password,
		acc.UserType,
		acc.UpdatedOn,
		id)
	if err != nil {
		return err
	}

	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from accounts")
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

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	sql := `select * from accounts where id=$1`
	rows, err := s.db.Query(sql, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account not found")
}

func (s *PostgresStore) FindAccountByEmail(email string) (*Account, error) {
	sql := `select * from accounts where email=$1`
	rows, err := s.db.Query(sql, email)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %s doesn't exist", email)
}

func (s *PostgresStore) UpdateAllTokens(token string, refreshToken string, id int) error {
	sqlGet := `select * from accounts where ID=$1`
	rows, err := s.db.Query(sqlGet, id)
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

	sqlUpdate := `update accounts set token=$1, refresh_token=$2 where id=$3`
	_, err = s.db.Exec(sqlUpdate, acc.Token, acc.RefreshToken, id)
	if err != nil {
		return err
	}

	return nil
}
func (s *PostgresStore) FindAccountByUuid(uuid string) (*Account, error) {
	sql := `select * from accounts where uuid=$1`
	rows, err := s.db.Query(sql, uuid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %s doesn't exist", uuid)
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
		&acc.RefreshToken,
		&acc.CreatedOn,
		&acc.UpdatedOn,
	)
	return acc, err
}
