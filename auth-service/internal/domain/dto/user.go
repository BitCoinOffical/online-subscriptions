package dto

type UsersRegisterDTO struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	Password_retry string `json:"password_retry"`
}
type UsersLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
