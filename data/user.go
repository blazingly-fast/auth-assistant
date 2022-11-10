package data

import ()

// func UpdateAllTokens(token string, refreshToken string, id int) {
// 	rows, err := DB.Query("select * from users where ID=$1", id)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "something went wrong updating tokens: %v\n", err)
// 		os.Exit(1)
// 	}
// 	user := User{}
// 	for rows.Next() {
// 		rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password)
// 	}
// 	user.Token = token
// 	user.Refresh_token = refreshToken

// 	sql := fmt.Sprintf("update users set token='%s', refresh_token='%s' where id=%d", user.Token, user.Refresh_token, id)

// 	DB.Exec(sql)
// }

// func Signup(u *User) {
// 	createSql := `
// 	  create table if not exists users(
// 	  id SERIAL PRIMARY KEY,
// 	  name text,
// 	  password text,
// 	  token text,
// 	  refresh_token text,
// 	  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
// 	  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
// 	  );
// 	  `
// 	_, err := DB.Exec(createSql)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
// 		os.Exit(1)

// 	}

// 	_, errz := DB.Exec("insert into users(name, password, token, refresh_token) values ($1, $2, $3, $4)", u.Name, u.Password, u.Token, u.Refresh_token)
// 	if errz != nil {
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "insetion failed: %v\n", errz)
// 			os.Exit(1)
// 		}
// 	}
// }

func login() {}
