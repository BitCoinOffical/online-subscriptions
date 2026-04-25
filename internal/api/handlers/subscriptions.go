package handlers

import (
	"net/http"
	"strconv"

	"github.com/BitCoinOffical/online-subscriptions/internal/api/response"
	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/online-subscriptions/internal/interfaces/http/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	service *services.SubscriptionService
	logger  *zap.Logger
}

func NewSubscriptionHandler(service *services.SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, logger: logger}
}

// CRUDL
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var dto dto.SubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}
	h.logger.Debug("received data", zap.Any("data", dto))

	if err := h.service.CreateSubscription(c.Request.Context(), &dto); err != nil {
		response.InternalServerError(c, err, "failed to create subscription", h.logger)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *SubscriptionHandler) GetSubscriptionsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}
	h.logger.Debug("received id", zap.Any("id", id))

	sub, err := h.service.GetSubscriptionsById(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err, "failed to get subscription by id", h.logger)
		return
	}
	h.logger.Debug("received data from data base", zap.Any("data", sub))

	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) UpdateSubscriptions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}
	h.logger.Debug("received id", zap.Any("id", id))

	var dto dto.SubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}
	h.logger.Debug("received data", zap.Any("data", dto))

	if err := h.service.UpdateSubscriptions(c.Request.Context(), &dto, id); err != nil {
		response.InternalServerError(c, err, "failed to update subscription by id", h.logger)
		return
	}

	c.Status(http.StatusOK)

}

func (h *SubscriptionHandler) DeleteSubscriptions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}
	h.logger.Debug("received id", zap.Any("id", id))

	if err := h.service.DeleteSubscriptions(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err, "failed delete subscriptions", h.logger)
		return
	}

	c.Status(http.StatusOK)
}

func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	subs, err := h.service.GetSubscriptions(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err, "failed to get subscriptions", h.logger)
		return
	}
	h.logger.Debug("received data from data base", zap.Any("data", subs))

	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) GetSubscriptionsFilter(c *gin.Context) {
	user_id := c.Query("user_id")
	service_name := c.Query("service_name")
	from := c.Query("from")
	to := c.Query("to")
	h.logger.Debug("received query", zap.Any("user_id", user_id), zap.Any("service_name", service_name), zap.Any("from", from), zap.Any("to", to))

	total, err := h.service.GetSubscriptionsFilter(c.Request.Context(), from, to, user_id, service_name)
	if err != nil {
		response.InternalServerError(c, err, "fail get total price", h.logger)
		return
	}
	h.logger.Debug("total price", zap.Int("total", total))

	c.JSON(http.StatusOK, total)
}
