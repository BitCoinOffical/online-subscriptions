package handlers

import (
	"net/http"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/api/response"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/domain/dto"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/interfaces/http/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type User struct {
	usersrvc *services.UserService
	logger   *zap.Logger
}

func NewUserHandler(usersrvc *services.UserService, logger *zap.Logger) *User {
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
		response.InternalServerError(c, err, "user register failed", h.logger)
		return
	}

	c.JSON(http.StatusCreated, tokens)
}
