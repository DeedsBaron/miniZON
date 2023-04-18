package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.SugaredLogger

func Init(devel bool, level string) {
	globalLogger = New(devel, level)
}

func New(devel bool, level string) *zap.SugaredLogger {
	var logger *zap.Logger
	var err error
	if devel {
		logger, err = zap.NewDevelopment()
	} else {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/2 15:04:05")

		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.EncoderConfig = encoderCfg

		switch level {
		case "DEBUG":
			cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		case "WARN":
			cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		case "ERROR":
			cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		case "INFO":
			cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		}

		logger, err = cfg.Build()
	}
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}

func Info(args ...interface{}) {
	globalLogger.Info(args)
}

func Debugf(msg string, args ...interface{}) {
	globalLogger.Debugf(msg, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	globalLogger.Debugw(msg, keysAndValues...)
}

func Infof(msg string, args ...interface{}) {
	globalLogger.Infof(msg, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	globalLogger.Infow(msg, keysAndValues...)
}

func Errorf(msg string, args ...interface{}) {
	globalLogger.Errorf(msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	globalLogger.Warnf(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	globalLogger.Fatalf(msg, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	globalLogger.Errorw(msg, keysAndValues...)
}
