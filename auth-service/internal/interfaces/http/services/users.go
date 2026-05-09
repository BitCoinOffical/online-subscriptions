package services

import (
	"context"
	"fmt"
	"time"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/dto"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/models"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/cache"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/repo"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/pkg/jwt"
	"github.com/google/uuid"
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

	if err := s.session.SaveToken(ctx, id, refreshToken, RefreshTTL); err != nil {
		return nil, fmt.Errorf("s.session.SaveToken: %w", err)
	}

	return &models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
func (s *UserService) LoginUser(ctx context.Context, user dto.UsersLoginDTO) (*models.Tokens, error) {
	userProfile, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, fmt.Errorf("s.userRepo.GetUserByEmail: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userProfile.Password_hash), []byte(user.Password))
	if err != nil {
		return nil, fmt.Errorf("bcrypt.CompareHashAndPassword: %w", domain.ErrInvalidCredentials)
	}

	accessToken, err := s.tokens.GenerateToken(userProfile.Id, AccessTTL)
	if err != nil {
		return nil, fmt.Errorf("accessToken s.tokens.GenerateToken: %w", err)
	}
	refreshToken, err := s.tokens.GenerateToken(userProfile.Id, RefreshTTL)
	if err != nil {
		return nil, fmt.Errorf("refreshToken s.tokens.GenerateToken: %w", err)
	}

	if err := s.session.SaveToken(ctx, userProfile.Id, refreshToken, RefreshTTL); err != nil {
		return nil, fmt.Errorf("s.session.SaveToken: %w", err)
	}

	return &models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) UpdateAccessToken(ctx context.Context, tokens dto.TokensDTO) (*models.Tokens, error) {
	refreshToken := tokens.RefreshToken

	user, err := s.tokens.ParseToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("s.tokens.ParseToken: %w", err)
	}

	userID, err := uuid.Parse(user.UserID)
	if err != nil {
		return nil, fmt.Errorf("uuid.Parse: %w", err)
	}

	savedToken, err := s.session.GetToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("s.session.GetToken: %w", err)
	}
	if savedToken != refreshToken {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := s.tokens.GenerateToken(userID, AccessTTL)
	if err != nil {
		return nil, fmt.Errorf("accessToken s.tokens.GenerateToken: %w", err)
	}
	return &models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (s *UserService) Logout(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("uuid.Parse: %w", err)
	}

	if err := s.session.DeleteRefreshToken(ctx, userID); err != nil {
		return fmt.Errorf("s.session.DeleteRefreshToken: %w", err)
	}

	return nil
}
