package data

import (
	"fmt"
	"os"
	"time"
)

var DB = Init()

type User struct {
	ID            int       `json:"id"`
	Name          string    `json:"name" validate:"required"`
	Password      string    `json:"password" validate:"required"`
	Token         string    `json:"token"`
	Refresh_token string    `json:"refresh_token"`
	CreatedOn     time.Time `json:"created_at"`
	UpdatedOn     time.Time `json:"updated_at"`
	DeletedOn     time.Time `json:"deleted_at"`
}

func UpdateAllTokens(token string, refreshToken string, id int) {
	rows, err := DB.Query("select * from users where ID=$1", id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "something went wrong updating tokens: %v\n", err)
		os.Exit(1)
	}
	user := User{}
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Password)
	}
	user.Token = token
	user.Refresh_token = refreshToken

	sql := fmt.Sprintf("update users set token='%s', refresh_token='%s' where id=%d", user.Token, user.Password, id)

	DB.Exec(sql)
}

func CheckEmail(u *User) *User {
	// check if user exist and store it in found user
	rows, err := DB.Query("select * from users where name=$1", u.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "something went wrong: %v\n", err)
		os.Exit(1)
	}
	user := User{}
	for rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Password)
	}
	return &user
}

func Signup(u *User) {
	createSql := `
	  create table if not exists users(
	  id SERIAL PRIMARY KEY,
	  name text,
	  password text,
	  token text,
	  refresh_token text,
	  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),	
	  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),	
	  deleted_at TIMESTAMPTZ
	  );
	  `
	_, err := DB.Exec(createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)

	}

	_, errz := DB.Exec("insert into users(name, password, token, refresh_token) values ($1, $2, $3, $4)", u.Name, u.Password, u.Token, u.Refresh_token)
	if errz != nil {
		if err != nil {
			fmt.Fprintf(os.Stderr, "insetion failed: %v\n", errz)
			os.Exit(1)
		}
	}
}

func login() {}
