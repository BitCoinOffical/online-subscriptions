package dto

// UsersRegisterDTO represents user registration request
type UsersRegisterDTO struct {
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required,min=8"`
	Password_retry string `json:"password_retry" binding:"required,min=8"`
}

// UsersLoginDTO represents user login request
type UsersLoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}
