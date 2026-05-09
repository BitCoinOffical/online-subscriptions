package api

import (
	"context"
	"net/http"

	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/api/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine *gin.Engine
	h      *handlers.Handlers
	srv    *http.Server
}

func NewServer(h *handlers.Handlers, port string) *Server {
	engine := gin.New()
	return &Server{
		h:      h,
		engine: engine,
		srv: &http.Server{
			Addr:    ":8080",
			Handler: engine,
		},
	}
}

func (s *Server) Run() error {
	s.engine.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/api/v1")
	})
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	subs := s.engine.Group("/api/v1")
	{
		subs.POST("/subscriptions", s.h.Subs.CreateSubscription)             //Create subscription
		subs.GET("/subscriptions/:id", s.h.Subs.GetSubscriptionsById)        //Get subscription
		subs.GET("/subscriptions/", s.h.Subs.GetSubscriptions)               //Get all subscriptions
		subs.GET("/subscriptions", s.h.Subs.GetSubscriptionsFilter)          //Get subscriptions with filter
		subs.PUT("/subscriptions/:id", s.h.Subs.FullUpdateSubscriptionsById) //Full update subscription
		subs.PATCH("/subscriptions/:id", s.h.Subs.UpdateSubscriptionsById)   //Update subscription
		subs.DELETE("/subscriptions/:id", s.h.Subs.DeleteSubscriptions)      //Delete subscription
	}

	return s.srv.ListenAndServe()
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
