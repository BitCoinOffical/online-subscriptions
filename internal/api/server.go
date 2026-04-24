package api

import (
	"net/http"

	"github.com/BitCoinOffical/online-subscriptions/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	h      *handlers.Handlers
}

func NewServer(h *handlers.Handlers) *Server {
	return &Server{h: h, engine: gin.New()}
}

func (s *Server) Run() error {
	s.engine.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/api/v1")
	})

	subs := s.engine.Group("/api/v1")
	{
		subs.POST("/subscriptions", s.h.Subs.CreateSubscription)        //Create subscription
		subs.GET("/subscriptions/:id", s.h.Subs.GetSubscriptionsById)   //Get subscription
		subs.GET("/subscriptions", s.h.Subs.GetSubscriptions)           //Get subscriptions
		subs.PATCH("/subscriptions/:id", s.h.Subs.UpdateSubscriptions)  //Update subscription
		subs.DELETE("/subscriptions/:id", s.h.Subs.DeleteSubscriptions) //Delete subscription
	}

	return s.engine.Run()
}
