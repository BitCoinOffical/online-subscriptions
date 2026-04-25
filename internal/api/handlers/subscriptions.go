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

// @Summary Create subscription
// @Tags subscriptions
// @Description Creates a subscription
// @Accept json
// @Produce json
// @Param request body dto.SubscriptionDTO true "Subscription data"
// @Success 201 {string} string "Created"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
	h.logger.Debug("received data", zap.Any("data", dto))

	if err := h.service.UpdateSubscriptionsById(c.Request.Context(), &dto, id); err != nil {
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
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscriptions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, err, "invalid subscription id", h.logger)
		return
	}
	h.logger.Debug("received id", zap.Any("id", id))

	if err := h.service.DeleteSubscriptions(c.Request.Context(), id); err != nil {
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
// @Failure 500 {object} ErrorResponse
// @Router /subscriptions/ [get]
func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {
	subs, err := h.service.GetSubscriptions(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err, "failed to get subscriptions", h.logger)
		return
	}
	h.logger.Debug("received data from data base", zap.Any("data", subs))

	c.JSON(http.StatusOK, subs)
}

// @Summary Calculating the total cost of subscriptions with filtering
// @Tags subscriptions
// @Description Calculating the total cost of subscriptions with filtering
// @Param user_id query string false "User UUID"
// @Param service_name query string false "Subscription service name"
// @Param from query string true "Start period (MM-YYYY) or (DD-MM-YYYY)"
// @Param to query string true "End period (MM-YYYY) or (DD-MM-YYYY)"
// @Produce json
// @Success 200 {integer} integer
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
