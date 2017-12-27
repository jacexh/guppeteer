package guppeteer

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

func init() {
	cfg := zap.NewProductionConfig()
	//cfg.Encoding = "console"
	cfg.Sampling = nil
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	Logger, _ = cfg.Build()
}
