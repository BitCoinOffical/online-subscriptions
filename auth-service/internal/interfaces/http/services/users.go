package services

import (
	"context"
	"fmt"
	"time"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/dto"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/models"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/cache"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/repo"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTTL  = 15 * time.Minute
	RefreshTTL = 24 * time.Hour
)

type UserService struct {
	userRepo *repo.UserRepo
	tokens   *jwt.ManagerToken
	session  *cache.Cache
}

func NewUserService(userRepo *repo.UserRepo, tokens *jwt.ManagerToken, session *cache.Cache) *UserService {
	return &UserService{userRepo: userRepo, tokens: tokens, session: session}
}

func (s *UserService) RegisterUser(ctx context.Context, user dto.UsersRegisterDTO) (*models.Tokens, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	userModel := &models.Users{
		Email:         user.Email,
		Password_hash: string(hashPass),
	}

	id, err := s.userRepo.RegisterUser(ctx, userModel)
	if err != nil {
		return nil, fmt.Errorf("s.userRepo.RegisterUser: %w", err)
	}

	accessToken, err := s.tokens.GenerateToken(id, AccessTTL)
	if err != nil {
		return nil, fmt.Errorf("accessToken s.tokens.GenerateToken: %w", err)
	}
	refreshToken, err := s.tokens.GenerateToken(id, RefreshTTL)
	if err != nil {
		return nil, fmt.Errorf("refreshToken s.tokens.GenerateToken: %w", err)
	}

	if err := s.session.SaveToken(ctx, id, refreshToken); err != nil {
		return nil, fmt.Errorf("s.session.SaveToken: %w", err)
	}

	return &models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
