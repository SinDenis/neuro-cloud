package domain

type User struct {
	Id       int64  `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
