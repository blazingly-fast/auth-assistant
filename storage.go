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

type Getter interface {
	GetAccounts() ([]*Account, error)
	GetAccountByField(string, any) (*Account, error)
}

type Putter interface {
	UpdateAccount(*UpdateAccountRequest, string) error
	UpdateAllTokens(string, string, int) error
}

type Deleter interface {
	DeleteAccount(string) error
}

type Poster interface {
	CreateAccout(*Account) error
}

type Storer interface {
	Getter
	Putter
	Deleter
	Poster
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
