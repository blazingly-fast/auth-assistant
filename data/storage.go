package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Getter interface {
	GetAccounts(int, int) (*AccountList, error)
	GetAccountByField(string, any) (*Account, error)
}

type Putter interface {
	UpdateAccount(*UpdateAccountRequest, string) error
	UpdateAllTokens(string, string, int) error
	UpdateAvatar(string, string) error
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

	postgresqlDbInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_SSLMODE"))

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
	  create table if not exists account(
	  id SERIAL PRIMARY KEY,
	  first_name text,
	  last_name text,
	  password text,
      email text,	
	  user_type text,
	  avatar text,
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
