package dto

type UsersRegisterDTO struct {
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required,min=8"`
	Password_retry string `json:"password_retry" binding:"required,min=8"`
}
type UsersLoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}
