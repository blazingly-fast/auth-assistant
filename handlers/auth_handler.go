package handlers

import "golang.org/x/crypto/bcrypt"

func (u *Users) HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		u.l.Panic(err)
	}
	return string(bytes)
}
