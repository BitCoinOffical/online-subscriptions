package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/dto"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/models"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/services"
	repo_mocks "github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/services/mock"
	jwtpkg "github.com/BitCoinOffical/online-subscriptions/auth-service/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

type testDeps struct {
	svc    *services.UserService
	repo   *repo_mocks.MockUserRepo
	cache  *repo_mocks.MockCache
	tokens *repo_mocks.MockManagerToken
}

func newDeps(t *testing.T) testDeps {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := repo_mocks.NewMockUserRepo(ctrl)
	cache := repo_mocks.NewMockCache(ctrl)
	tokens := repo_mocks.NewMockManagerToken(ctrl)
	svc := services.NewUserService(repo, tokens, cache)
	return testDeps{svc: svc, repo: repo, cache: cache, tokens: tokens}
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)
	return string(hash)
}

func TestRegisterUser_Success(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()

	d.repo.EXPECT().
		RegisterUser(gomock.Any(), gomock.Any()).
		Return(userID, nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.AccessTTL).
		Return("access-token", nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.RefreshTTL).
		Return("refresh-token", nil)

	d.cache.EXPECT().
		SaveToken(gomock.Any(), userID, "refresh-token", services.RefreshTTL).
		Return(nil)

	tokens, err := d.svc.RegisterUser(context.Background(), dto.UsersRegisterDTO{
		Email:    "test@example.com",
		Password: "secret",
	})

	require.NoError(t, err)
	assert.Equal(t, "access-token", tokens.AccessToken)
	assert.Equal(t, "refresh-token", tokens.RefreshToken)
}

func TestRegisterUser_RepoError(t *testing.T) {
	d := newDeps(t)
	dbErr := errors.New("db error")

	d.repo.EXPECT().
		RegisterUser(gomock.Any(), gomock.Any()).
		Return(uuid.UUID{}, dbErr)

	tokens, err := d.svc.RegisterUser(context.Background(), dto.UsersRegisterDTO{
		Email:    "test@example.com",
		Password: "secret",
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "s.userRepo.RegisterUser")
	assert.ErrorIs(t, err, dbErr)
}

func TestRegisterUser_GenerateAccessTokenError(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	tokenErr := errors.New("token error")

	d.repo.EXPECT().
		RegisterUser(gomock.Any(), gomock.Any()).
		Return(userID, nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.AccessTTL).
		Return("", tokenErr)

	tokens, err := d.svc.RegisterUser(context.Background(), dto.UsersRegisterDTO{
		Email:    "test@example.com",
		Password: "secret",
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "accessToken s.tokens.GenerateToken")
	assert.ErrorIs(t, err, tokenErr)
}

func TestRegisterUser_GenerateRefreshTokenError(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	tokenErr := errors.New("token error")

	d.repo.EXPECT().
		RegisterUser(gomock.Any(), gomock.Any()).
		Return(userID, nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.AccessTTL).
		Return("access-token", nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.RefreshTTL).
		Return("", tokenErr)

	tokens, err := d.svc.RegisterUser(context.Background(), dto.UsersRegisterDTO{
		Email:    "test@example.com",
		Password: "secret",
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "refreshToken s.tokens.GenerateToken")
	assert.ErrorIs(t, err, tokenErr)
}

func TestRegisterUser_SaveTokenError(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	cacheErr := errors.New("cache error")

	d.repo.EXPECT().
		RegisterUser(gomock.Any(), gomock.Any()).
		Return(userID, nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.AccessTTL).
		Return("access-token", nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.RefreshTTL).
		Return("refresh-token", nil)

	d.cache.EXPECT().
		SaveToken(gomock.Any(), userID, "refresh-token", services.RefreshTTL).
		Return(cacheErr)

	tokens, err := d.svc.RegisterUser(context.Background(), dto.UsersRegisterDTO{
		Email:    "test@example.com",
		Password: "secret",
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "s.session.SaveToken")
	assert.ErrorIs(t, err, cacheErr)
}

func TestLoginUser_Success(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	password := "secret"

	d.repo.EXPECT().
		GetUserByEmail(gomock.Any(), "test@example.com").
		Return(&models.Users{
			Id:            userID,
			Email:         "test@example.com",
			Password_hash: hashPassword(t, password),
		}, nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.AccessTTL).
		Return("access-token", nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.RefreshTTL).
		Return("refresh-token", nil)

	d.cache.EXPECT().
		SaveToken(gomock.Any(), userID, "refresh-token", services.RefreshTTL).
		Return(nil)

	tokens, err := d.svc.LoginUser(context.Background(), dto.UsersLoginDTO{
		Email:    "test@example.com",
		Password: password,
	})

	require.NoError(t, err)
	assert.Equal(t, "access-token", tokens.AccessToken)
	assert.Equal(t, "refresh-token", tokens.RefreshToken)
}

func TestLoginUser_GetUserByEmailError(t *testing.T) {
	d := newDeps(t)
	dbErr := errors.New("db error")

	d.repo.EXPECT().
		GetUserByEmail(gomock.Any(), "test@example.com").
		Return(nil, dbErr)

	tokens, err := d.svc.LoginUser(context.Background(), dto.UsersLoginDTO{
		Email:    "test@example.com",
		Password: "secret",
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "s.userRepo.GetUserByEmail")
	assert.ErrorIs(t, err, dbErr)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()

	d.repo.EXPECT().
		GetUserByEmail(gomock.Any(), "test@example.com").
		Return(&models.Users{
			Id:            userID,
			Email:         "test@example.com",
			Password_hash: hashPassword(t, "correct"),
		}, nil)

	tokens, err := d.svc.LoginUser(context.Background(), dto.UsersLoginDTO{
		Email:    "test@example.com",
		Password: "wrong",
	})

	assert.Nil(t, tokens)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLoginUser_GenerateAccessTokenError(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	password := "secret"
	tokenErr := errors.New("token error")

	d.repo.EXPECT().
		GetUserByEmail(gomock.Any(), "test@example.com").
		Return(&models.Users{
			Id:            userID,
			Email:         "test@example.com",
			Password_hash: hashPassword(t, password),
		}, nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.AccessTTL).
		Return("", tokenErr)

	tokens, err := d.svc.LoginUser(context.Background(), dto.UsersLoginDTO{
		Email:    "test@example.com",
		Password: password,
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "accessToken s.tokens.GenerateToken")
	assert.ErrorIs(t, err, tokenErr)
}

func TestLoginUser_GenerateRefreshTokenError(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	password := "secret"
	tokenErr := errors.New("token error")

	d.repo.EXPECT().
		GetUserByEmail(gomock.Any(), "test@example.com").
		Return(&models.Users{
			Id:            userID,
			Email:         "test@example.com",
			Password_hash: hashPassword(t, password),
		}, nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.AccessTTL).
		Return("access-token", nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.RefreshTTL).
		Return("", tokenErr)

	tokens, err := d.svc.LoginUser(context.Background(), dto.UsersLoginDTO{
		Email:    "test@example.com",
		Password: password,
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "refreshToken s.tokens.GenerateToken")
	assert.ErrorIs(t, err, tokenErr)
}

func TestLoginUser_SaveTokenError(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	password := "secret"
	cacheErr := errors.New("cache error")

	d.repo.EXPECT().
		GetUserByEmail(gomock.Any(), "test@example.com").
		Return(&models.Users{
			Id:            userID,
			Email:         "test@example.com",
			Password_hash: hashPassword(t, password),
		}, nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.AccessTTL).
		Return("access-token", nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.RefreshTTL).
		Return("refresh-token", nil)

	d.cache.EXPECT().
		SaveToken(gomock.Any(), userID, "refresh-token", services.RefreshTTL).
		Return(cacheErr)

	tokens, err := d.svc.LoginUser(context.Background(), dto.UsersLoginDTO{
		Email:    "test@example.com",
		Password: password,
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "s.session.SaveToken")
	assert.ErrorIs(t, err, cacheErr)
}

func TestUpdateAccessToken_Success(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	refreshToken := "valid-refresh-token"

	d.tokens.EXPECT().
		ParseToken(refreshToken).
		Return(&jwtpkg.Claims{UserID: userID.String()}, nil)

	d.cache.EXPECT().
		GetToken(gomock.Any(), userID).
		Return(refreshToken, nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.AccessTTL).
		Return("new-access-token", nil)

	tokens, err := d.svc.UpdateAccessToken(context.Background(), dto.TokensDTO{
		RefreshToken: refreshToken,
	})

	require.NoError(t, err)
	assert.Equal(t, "new-access-token", tokens.AccessToken)
	assert.Equal(t, refreshToken, tokens.RefreshToken)
}

func TestUpdateAccessToken_ParseTokenError(t *testing.T) {
	d := newDeps(t)
	parseErr := errors.New("invalid token")

	d.tokens.EXPECT().
		ParseToken("bad-token").
		Return(nil, parseErr)

	tokens, err := d.svc.UpdateAccessToken(context.Background(), dto.TokensDTO{
		RefreshToken: "bad-token",
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "s.tokens.ParseToken")
	assert.ErrorIs(t, err, parseErr)
}

func TestUpdateAccessToken_InvalidUserIDInClaims(t *testing.T) {
	d := newDeps(t)
	refreshToken := "valid-token"

	d.tokens.EXPECT().
		ParseToken(refreshToken).
		Return(&jwtpkg.Claims{UserID: "not-a-uuid"}, nil)

	tokens, err := d.svc.UpdateAccessToken(context.Background(), dto.TokensDTO{
		RefreshToken: refreshToken,
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "uuid.Parse")
}

func TestUpdateAccessToken_GetTokenError(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	refreshToken := "valid-refresh-token"
	cacheErr := errors.New("cache error")

	d.tokens.EXPECT().
		ParseToken(refreshToken).
		Return(&jwtpkg.Claims{UserID: userID.String()}, nil)

	d.cache.EXPECT().
		GetToken(gomock.Any(), userID).
		Return("", cacheErr)

	tokens, err := d.svc.UpdateAccessToken(context.Background(), dto.TokensDTO{
		RefreshToken: refreshToken,
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "s.session.GetToken")
	assert.ErrorIs(t, err, cacheErr)
}

func TestUpdateAccessToken_TokenMismatch(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	refreshToken := "valid-refresh-token"

	d.tokens.EXPECT().
		ParseToken(refreshToken).
		Return(&jwtpkg.Claims{UserID: userID.String()}, nil)

	d.cache.EXPECT().
		GetToken(gomock.Any(), userID).
		Return("different-token", nil)

	tokens, err := d.svc.UpdateAccessToken(context.Background(), dto.TokensDTO{
		RefreshToken: refreshToken,
	})

	assert.Nil(t, tokens)
	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestUpdateAccessToken_GenerateAccessTokenError(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	refreshToken := "valid-refresh-token"
	tokenErr := errors.New("token error")

	d.tokens.EXPECT().
		ParseToken(refreshToken).
		Return(&jwtpkg.Claims{UserID: userID.String()}, nil)

	d.cache.EXPECT().
		GetToken(gomock.Any(), userID).
		Return(refreshToken, nil)

	d.tokens.EXPECT().
		GenerateToken(userID, services.AccessTTL).
		Return("", tokenErr)

	tokens, err := d.svc.UpdateAccessToken(context.Background(), dto.TokensDTO{
		RefreshToken: refreshToken,
	})

	assert.Nil(t, tokens)
	assert.ErrorContains(t, err, "accessToken s.tokens.GenerateToken")
	assert.ErrorIs(t, err, tokenErr)
}

func TestLogout_Success(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()

	d.cache.EXPECT().
		DeleteRefreshToken(gomock.Any(), userID).
		Return(nil)

	err := d.svc.Logout(context.Background(), userID.String())

	assert.NoError(t, err)
}

func TestLogout_InvalidUUID(t *testing.T) {
	d := newDeps(t)

	err := d.svc.Logout(context.Background(), "not-a-uuid")

	assert.ErrorContains(t, err, "uuid.Parse")
}

func TestLogout_DeleteTokenError(t *testing.T) {
	d := newDeps(t)
	userID := uuid.New()
	cacheErr := errors.New("cache error")

	d.cache.EXPECT().
		DeleteRefreshToken(gomock.Any(), userID).
		Return(cacheErr)

	err := d.svc.Logout(context.Background(), userID.String())

	assert.ErrorContains(t, err, "s.session.DeleteRefreshToken")
	assert.ErrorIs(t, err, cacheErr)
}
