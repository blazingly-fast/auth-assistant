package data

type User struct {
	ID            int    `json:"id"`
	Name          string `json:"name" validate:"required"`
	Password      string `json:"password" validate:"required"`
	Token         string `json:"token"`
	Refresh_token string `json:"refresh_token"`
	CreatedOn     string `json:"-"`
	UpdatedOn     string `json:"-"`
	DeletedOn     string `json:"-"`
}

func Signup(u *User) {

}

func login() {}
