package api

import (
	"context"
	"net/http"
	"time"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/api/handlers"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/api/middleware"
	"github.com/BitCoinOffical/online-subscriptions/auth-service/pkg/jwt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

const (
	HeaderTimeout = 5
)

type Server struct {
	engine   *gin.Engine
	h        *handlers.Handlers
	srv      *http.Server
	manager  *jwt.ManagerToken
	logger   *zap.Logger
	limitter *middleware.RateLimiter
}

func NewServer(h *handlers.Handlers, manager *jwt.ManagerToken, limitter *middleware.RateLimiter, port string, logger *zap.Logger) *Server {
	engine := gin.New()
	return &Server{
		h:        h,
		engine:   engine,
		manager:  manager,
		limitter: limitter,
		logger:   logger,
		srv: &http.Server{
			Addr:              ":" + port,
			Handler:           engine,
			ReadHeaderTimeout: HeaderTimeout * time.Second,
		},
	}
}

func (s *Server) Run() error {
	s.engine.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/api/v1")
	})
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := s.engine.Group("/auth")
	auth.Use(s.limitter.RateLimiter())
	{
		auth.POST("/register", s.h.User.RegisterUser)
		auth.POST("/login", s.h.User.LoginUser)
		auth.POST("/refresh", s.h.User.UpdateAccessToken)
		auth.DELETE("/logout", middleware.AuthMiddleware(s.manager, s.logger), s.h.User.Logout)
	}

	return s.srv.ListenAndServe()
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
