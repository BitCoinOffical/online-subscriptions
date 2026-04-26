package zaplogger

import (
	"errors"

	"github.com/BitCoinOffical/online-subscriptions/config"
	"go.uber.org/zap"
)

const (
	EnvProd = "Prod"
	EnvDev  = "Dev"
)

func NewLogger(cfg *config.AppConfig) (*zap.Logger, error) {
	switch cfg.DebugLevel {
	case "Dev":
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		return logger, nil
	case "Prod":
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		return logger, nil
	default:
		return nil, errors.New("incorrect debug value")
	}
}
