package handlers

import (
	"errors"
	"net/http"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/api/response"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type User struct {
	usersrvc UserService
	logger   *zap.Logger
}

func NewUserHandler(usersrvc UserService, logger *zap.Logger) *User {
	return &User{logger: logger, usersrvc: usersrvc}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UsersRegisterDTO true "User registration data"
// @Success 201 {object} dto.TokensDTO "Successfully registered, returns access and refresh tokens"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or password mismatch"
// @Failure 409 {object} response.ErrorResponse "Email already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /register [post]
func (h *User) RegisterUser(c *gin.Context) {
	var user dto.UsersRegisterDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, err, "user register failed", h.logger)
		return
	}

	if user.Password != user.Password_retry {
		response.BadRequest(c, domain.ErrPasswordMismatch, "pass error", h.logger)
		return
	}

	tokens, err := h.usersrvc.RegisterUser(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.Conflict(c, err, "email already exists", h.logger)
			return
		}
		response.InternalServerError(c, err, "user register failed", h.logger)
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

// LoginUser godoc
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.UsersLoginDTO true "User login credentials"
// @Success 200 {object} dto.TokensDTO "Successfully logged in, returns access and refresh tokens"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /login [post]
func (h *User) LoginUser(c *gin.Context) {
	var user dto.UsersLoginDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, err, "user login failed", h.logger)
		return
	}

	tokens, err := h.usersrvc.LoginUser(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, domain.ErrNotFound) {
			response.Unauthorized(c, err, "invalid credentials", h.logger)
			return
		}
		response.InternalServerError(c, err, "user login failed", h.logger)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// UpdateAccessToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.TokensDTO true "Refresh token"
// @Success 200 {object} dto.TokensDTO "Successfully refreshed, returns new access and refresh tokens"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired refresh token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /refresh [post]
func (h *User) UpdateAccessToken(c *gin.Context) {
	var token dto.TokensDTO
	if err := c.ShouldBindJSON(&token); err != nil {
		response.BadRequest(c, err, "user update access token failed", h.logger)
		return
	}

	tokens, err := h.usersrvc.UpdateAccessToken(c.Request.Context(), token)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, domain.ErrNotFound) {
			response.Unauthorized(c, err, "invalid credentials", h.logger)
			return
		}
		response.InternalServerError(c, err, "user update access token failed", h.logger)
		return
	}

	c.JSON(http.StatusOK, tokens)
}
// Logout godoc
// @Summary Logout user
// @Description Invalidate user session and tokens
// @Tags auth
// @Security BearerAuth
// @Success 204 "Successfully logged out"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - invalid or missing token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /logout [delete]
func (h *User) Logout(c *gin.Context) {
	id := c.GetString("user_id")

	if err := h.usersrvc.Logout(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err, "user logout failed", h.logger)
		return
	}

	c.Status(http.StatusNoContent)
}
