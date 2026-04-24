package handlers

import "github.com/gin-gonic/gin"

type SubscriptionHandler struct {
}

func NewSubscriptionHandler() *SubscriptionHandler {
	return &SubscriptionHandler{}
}

// CRUDL
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {

}

func (h *SubscriptionHandler) GetSubscriptionsById(c *gin.Context) {

}

func (h *SubscriptionHandler) UpdateSubscriptions(c *gin.Context) {

}

func (h *SubscriptionHandler) DeleteSubscriptions(c *gin.Context) {

}

func (h *SubscriptionHandler) GetSubscriptions(c *gin.Context) {

}
