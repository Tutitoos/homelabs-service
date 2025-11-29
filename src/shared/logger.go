package shared

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Zap    *zap.Logger
	Logger *zap.SugaredLogger
)

func InitLogger() {
	CapturePanic()

	var err error

	if Config.AppEnv == "production" {
		zapConfig := zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}
		Zap, err = zapConfig.Build()
	} else {
		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}
		Zap, err = zapConfig.Build()
	}

	if err != nil {
		panic(err)
	}
	defer Zap.Sync()

	Logger = Zap.Sugar()
	Logger.Info("Logger initialized successfully")
}
