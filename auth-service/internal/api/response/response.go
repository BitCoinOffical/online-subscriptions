package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func InternalServerError(c *gin.Context, err error, msg string, logger *zap.Logger) {
	logger.Error(msg, zap.Error(err), zap.String("path", c.FullPath()))
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: err.Error(),
	})
}

func BadRequest(c *gin.Context, err error, msg string, logger *zap.Logger) {
	logger.Info(msg, zap.Error(err), zap.String("path", c.FullPath()))
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: err.Error(),
	})
}

func NotFound(c *gin.Context, err error, msg string, logger *zap.Logger) {
	logger.Info(msg, zap.Error(err), zap.String("path", c.FullPath()))
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: err.Error(),
	})
}

func Unauthorized(c *gin.Context, err error, msg string, logger *zap.Logger) {
	logger.Info(msg, zap.Error(err), zap.String("path", c.FullPath()))
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: err.Error(),
	})
}

func Conflict(c *gin.Context, err error, msg string, logger *zap.Logger) {
	logger.Info(msg, zap.Error(err), zap.String("path", c.FullPath()))
	c.JSON(http.StatusConflict, ErrorResponse{
		Error: err.Error(),
	})
}

func ManyRequest(c *gin.Context, err error, msg string, logger *zap.Logger) {
	logger.Info(msg, zap.Error(err), zap.String("path", c.FullPath()))
	c.JSON(http.StatusTooManyRequests, ErrorResponse{
		Error: err.Error(),
	})
}
