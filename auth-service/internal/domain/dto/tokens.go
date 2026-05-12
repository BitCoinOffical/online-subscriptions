package dto

// TokensDTO represents JWT tokens response
type TokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
