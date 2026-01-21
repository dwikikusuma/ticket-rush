package domain

import "time"

type LoginRequest struct {
	Email    string `json:"email,required,email"`
	Password string `json:"password,required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	password  string
	CreatedAt time.Time `json:"createdAt"`
}

func (u *User) GSetPassword() string {
	return u.password
}
func (u *User) SetPassword(pw string) {
	u.password = pw
}
