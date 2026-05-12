package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/api/response"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	service SubscriptionService
	logger  *zap.Logger
}

func NewSubscriptionHandler(service SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, logger: logger}
}

// @Summary Create subscription
// @Tags subscriptions
// @Description Creates a subscription
// @Accept json
// @Produce json
// @Param request body dto.SubscriptionDTO true "Subscription data"
// @Success 201 {string} string "Created"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions [post]
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

// @Summary Get subscription by ID
// @Tags subscriptions
// @Description Get subscription by ID
// @Param id path int true "Subscription ID"
// @Produce json
// @Success 200 {object} models.Subscription
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscriptionsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid subscription id", h.logger)
		return
	}
	h.logger.Debug("received id", zap.Any("id", id))

	sub, err := h.service.GetSubscriptionsById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, err, "not found", h.logger)
			return
		}
		response.InternalServerError(c, err, "failed to get subscription by id", h.logger)
		return
	}
	h.logger.Debug("received data from data base", zap.Any("data", sub))

	c.JSON(http.StatusOK, sub)
}

// @Summary Partially update subscription
// @Tags subscriptions
// @Description Partially update subscription by ID
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param request body dto.SubscriptionDTO true "Fields to update"
// @Success 200 {string} string "Updated"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/{id} [patch]
func (h *SubscriptionHandler) UpdateSubscriptionsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid subscription id", h.logger)
		return
	}
	h.logger.Debug("received id", zap.Any("id", id))

	var dto dto.PatchSubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}
	if dto.IsEmpty() {
		response.BadRequest(c, domain.ErrEmptyPayload, "empty patch payload", h.logger)
		return
	}
	h.logger.Debug("received data", zap.Any("data", dto))

	if err := h.service.UpdateSubscriptionsById(c.Request.Context(), &dto, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, err, "not found", h.logger)
			return
		}
		response.InternalServerError(c, err, "failed to update subscription by id", h.logger)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Update subscription
// @Tags subscriptions
// @Description Update subscription by ID
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param request body dto.SubscriptionDTO true "Updated subscription data"
// @Success 200 {string} string "Updated"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) FullUpdateSubscriptionsById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid subscription id", h.logger)
		return
	}
	h.logger.Debug("received id", zap.Any("id", id))

	var dto dto.SubscriptionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, err, "invalid request body", h.logger)
		return
	}
	h.logger.Debug("received data", zap.Any("data", dto))

	if err := h.service.FullUpdateSubscriptionsById(c.Request.Context(), &dto, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, err, "not found", h.logger)
			return
		}
		response.InternalServerError(c, err, "failed to update subscription by id", h.logger)
		return
	}

	c.Status(http.StatusOK)

}

// @Summary Delete subscription
// @Tags subscriptions
// @Description Delete subscription by ID
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscriptions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid subscription id", h.logger)
		return
	}
	h.logger.Debug("received id", zap.Any("id", id))

	if err := h.service.DeleteSubscriptions(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(c, err, "not found", h.logger)
			return
		}
		response.InternalServerError(c, err, "failed delete subscriptions", h.logger)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Get all subscriptions
// @Tags subscriptions
// @Description Get all subscriptions
// @Produce json
// @Success 200 {array} models.Subscription
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/ [get]
func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		response.BadRequest(c, err, "invalid page", h.logger)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		response.BadRequest(c, err, "invalid limit", h.logger)
		return
	}

	offset := (page - 1) * limit

	subs, err := h.service.GetSubscriptions(c.Request.Context(), limit, offset)

	if err != nil {
		response.InternalServerError(c, err, "failed to get subscriptions", h.logger)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"data":  subs,
	})
}

// @Summary Calculating the total cost of subscriptions with filtering
// @Tags subscriptions
// @Description Calculating the total cost of subscriptions with filtering
// @Param user_id query string false "User UUID"
// @Param service_name query string false "Subscription service name"
// @Param from query string false "Start period (MM-YYYY) or (DD-MM-YYYY)"
// @Param to query string false "End period (MM-YYYY) or (DD-MM-YYYY)"
// @Produce json
// @Success 200 {integer} integer
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions [get]
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
