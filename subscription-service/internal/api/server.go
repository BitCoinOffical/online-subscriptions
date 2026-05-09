package api

import (
	"context"
	"net/http"

	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/api/handlers"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/internal/api/middleware"
	"github.com/BitCoinOffical/online-subscriptions/subscription-service/pkg/jwt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Server struct {
	engine   *gin.Engine
	h        *handlers.Handlers
	srv      *http.Server
	manager  *jwt.ManagerToken
	logger   *zap.Logger
	limitter *middleware.RateLimiter
}

func NewServer(h *handlers.Handlers, port string, manager *jwt.ManagerToken, limitter *middleware.RateLimiter, logger *zap.Logger) *Server {
	engine := gin.New()
	return &Server{
		h:        h,
		engine:   engine,
		manager:  manager,
		logger:   logger,
		limitter: limitter,
		srv: &http.Server{
			Addr:    ":" + port,
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
	subs.Use(middleware.AuthMiddleware(s.manager, s.logger), s.limitter.RateLimiter())
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
