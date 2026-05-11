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
func (h *User) Logout(c *gin.Context) {
	id := c.GetString("user_id")

	if err := h.usersrvc.Logout(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err, "user logout failed", h.logger)
		return
	}

	c.Status(http.StatusNoContent)
}
