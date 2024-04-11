package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var Logger *zap.Logger

func InitLogger(logMode string) {
	var (
		loggerInit *zap.Logger
		err        error
	)

	if logMode == "debug" {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		option := zap.AddStacktrace(zap.DPanicLevel)
		loggerInit, err = config.Build(option)
		if err != nil {
			log.Fatal("error creating logger : ", err.Error())
		}
		loggerInit.Debug("Logger started", zap.String("mode", "debug"))
	} else {
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		option := zap.AddStacktrace(zap.DPanicLevel)
		loggerInit, err = config.Build(option)
		if err != nil {
			log.Fatal("error creating logger : ", err.Error())
		}
		loggerInit.Info("Logger started", zap.String("mode", "production"))
	}

	Logger = loggerInit
}
