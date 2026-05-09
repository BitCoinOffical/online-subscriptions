package api

import (
	"context"
	"net/http"

	"github.com/BitCoinOffical/online-subscriptions/auth-service/internal/api/handlers"
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

	auth := s.engine.Group("/auth")
	{
		auth.POST("/register", s.h.User.RegisterUser)
		auth.POST("/login", s.h.User.RegisterUser)
		auth.GET("/logout", s.h.User.RegisterUser)
	}

	return s.srv.ListenAndServe()
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
