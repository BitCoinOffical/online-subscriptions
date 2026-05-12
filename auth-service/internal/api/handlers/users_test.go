package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/api/handlers"
	mocks "github.com/BitCoinOffical/online-subscriptions/auth-service/internal/api/handlers/mock"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/dto"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newRouter(h *handlers.User) *gin.Engine {
	r := gin.New()
	r.POST("/register", h.RegisterUser)
	r.POST("/login", h.LoginUser)
	r.POST("/token/refresh", h.UpdateAccessToken)
	r.POST("/logout", func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		h.Logout(c)
	})
	return r
}

func newTestHandler(t *testing.T) (*handlers.User, *mocks.MockUserService, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	svc := mocks.NewMockUserService(ctrl)
	h := handlers.NewUserHandler(svc, zap.NewNop())
	return h, svc, ctrl
}

func makeRequest(t *testing.T, router *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("failed to encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestRegisterUser_Success(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.UsersRegisterDTO{
		Email:          "test@example.com",
		Password:       "secret123",
		Password_retry: "secret123",
	}
	want := &models.Tokens{AccessToken: "access", RefreshToken: "refresh"}

	svc.EXPECT().
		RegisterUser(gomock.Any(), input).
		Return(want, nil)

	w := makeRequest(t, newRouter(h), http.MethodPost, "/register", input)

	assert.Equal(t, http.StatusCreated, w.Code)

	var got models.Tokens
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, *want, got)
}

func TestRegisterUser_PasswordMismatch(t *testing.T) {
	h, _, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.UsersRegisterDTO{
		Email:          "test@example.com",
		Password:       "secret",
		Password_retry: "other",
	}

	w := makeRequest(t, newRouter(h), http.MethodPost, "/register", input)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterUser_InvalidBody(t *testing.T) {
	h, _, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newRouter(h).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.UsersRegisterDTO{
		Email:          "taken@example.com",
		Password:       "secret123",
		Password_retry: "secret123",
	}

	svc.EXPECT().
		RegisterUser(gomock.Any(), input).
		Return(nil, domain.ErrEmailAlreadyExists)

	w := makeRequest(t, newRouter(h), http.MethodPost, "/register", input)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegisterUser_InternalError(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.UsersRegisterDTO{
		Email:          "test@example.com",
		Password:       "secret123",
		Password_retry: "secret123",
	}

	svc.EXPECT().
		RegisterUser(gomock.Any(), input).
		Return(nil, errors.New("db down"))

	w := makeRequest(t, newRouter(h), http.MethodPost, "/register", input)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLoginUser_Success(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.UsersLoginDTO{Email: "test@example.com", Password: "secret123"}
	want := &models.Tokens{AccessToken: "access", RefreshToken: "refresh"}

	svc.EXPECT().
		LoginUser(gomock.Any(), input).
		Return(want, nil)

	w := makeRequest(t, newRouter(h), http.MethodPost, "/login", input)

	assert.Equal(t, http.StatusOK, w.Code)

	var got models.Tokens
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, *want, got)
}

func TestLoginUser_InvalidBody(t *testing.T) {
	h, _, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newRouter(h).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.UsersLoginDTO{Email: "test@example.com", Password: "wrong123"}

	svc.EXPECT().
		LoginUser(gomock.Any(), input).
		Return(nil, domain.ErrInvalidCredentials)

	w := makeRequest(t, newRouter(h), http.MethodPost, "/login", input)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginUser_NotFound(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.UsersLoginDTO{Email: "ghost@example.com", Password: "secret123"}

	svc.EXPECT().
		LoginUser(gomock.Any(), input).
		Return(nil, domain.ErrNotFound)

	w := makeRequest(t, newRouter(h), http.MethodPost, "/login", input)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginUser_InternalError(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.UsersLoginDTO{Email: "test@example.com", Password: "secret123"}

	svc.EXPECT().
		LoginUser(gomock.Any(), input).
		Return(nil, errors.New("db down"))

	w := makeRequest(t, newRouter(h), http.MethodPost, "/login", input)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateAccessToken_Success(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.TokensDTO{RefreshToken: "old-refresh"}
	want := &models.Tokens{AccessToken: "new-access", RefreshToken: "new-refresh"}

	svc.EXPECT().
		UpdateAccessToken(gomock.Any(), input).
		Return(want, nil)

	w := makeRequest(t, newRouter(h), http.MethodPost, "/token/refresh", input)

	assert.Equal(t, http.StatusOK, w.Code)

	var got models.Tokens
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, *want, got)
}

func TestUpdateAccessToken_InvalidBody(t *testing.T) {
	h, _, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	req := httptest.NewRequest(http.MethodPost, "/token/refresh", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	newRouter(h).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAccessToken_InvalidCredentials(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.TokensDTO{RefreshToken: "expired"}

	svc.EXPECT().
		UpdateAccessToken(gomock.Any(), input).
		Return(nil, domain.ErrInvalidCredentials)

	w := makeRequest(t, newRouter(h), http.MethodPost, "/token/refresh", input)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateAccessToken_NotFound(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.TokensDTO{RefreshToken: "unknown"}

	svc.EXPECT().
		UpdateAccessToken(gomock.Any(), input).
		Return(nil, domain.ErrNotFound)

	w := makeRequest(t, newRouter(h), http.MethodPost, "/token/refresh", input)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateAccessToken_InternalError(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	input := dto.TokensDTO{RefreshToken: "token"}

	svc.EXPECT().
		UpdateAccessToken(gomock.Any(), input).
		Return(nil, errors.New("db down"))

	w := makeRequest(t, newRouter(h), http.MethodPost, "/token/refresh", input)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLogout_Success(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	svc.EXPECT().
		Logout(gomock.Any(), "test-user-id").
		Return(nil)

	w := makeRequest(t, newRouter(h), http.MethodPost, "/logout", nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestLogout_InternalError(t *testing.T) {
	h, svc, ctrl := newTestHandler(t)
	defer ctrl.Finish()

	svc.EXPECT().
		Logout(gomock.Any(), "test-user-id").
		Return(errors.New("db down"))

	w := makeRequest(t, newRouter(h), http.MethodPost, "/logout", nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
