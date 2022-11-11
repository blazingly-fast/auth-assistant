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
	insert into accounts(first_name, last_name, email, password, token, refresh_token)
	values($1, $2, $3, $4, $5, $6)
	`
	_, err := s.db.Exec(sql, acc.FirstName, acc.LastName, acc.Email, acc.Password, acc.Token, acc.RefreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	return nil, nil

}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	sql := fmt.Sprintf("select * from accounts where id=%d", id)
	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	acc := Account{}
	for rows.Next() {
		err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Password, &acc.Email, &acc.Token, &acc.RefreshToken, &acc.CreatedOn, &acc.UpdatedOn)
		if err != nil {
			return nil, err
		}
	}
	return &acc, nil
}

func (s *PostgresStore) CheckEmail(r *LoginRequest) (*Account, error) {
	// check if user exist and store it in found user
	sql := fmt.Sprintf("select * from accounts where email='%s'", r.Email)
	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	acc := Account{}
	for rows.Next() {
		err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Password, &acc.Email, &acc.Token, &acc.RefreshToken, &acc.CreatedOn, &acc.UpdatedOn)
		if err != nil {
			return nil, err
		}
	}
	return &acc, nil
}

func (s *PostgresStore) UpdateAllTokens(token string, refreshToken string, id int) error {
	rows, err := s.db.Query("select * from accounts where ID=$1", id)
	if err != nil {
		return err
	}
	acc := Account{}
	for rows.Next() {
		err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Password, &acc.Email, &acc.Token, &acc.RefreshToken, &acc.CreatedOn, &acc.UpdatedOn)
		if err != nil {
			return err
		}

	}
	acc.Token = token
	acc.RefreshToken = refreshToken

	sql := fmt.Sprintf("update users set token='%s', refresh_token='%s' where id=%d", acc.Token, acc.RefreshToken, id)

	_, err = s.db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
