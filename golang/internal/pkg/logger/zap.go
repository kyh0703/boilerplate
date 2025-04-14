package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap *zap.Logger = nil

func init() {
	var (
		encoding      string
		encoderConfig zapcore.EncoderConfig
	)
	if os.Getenv("env") == "prod" {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoding = "json"
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoding = "console"
	}
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	zapConfig := zap.Config{
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Level:             getLogLevel("debug"),
		Encoding:          encoding,
		EncoderConfig:     encoderConfig,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{},
	}
	z, err := zapConfig.Build()
	if err != nil {
		Zap = zap.NewExample()
	} else {
		Zap = z
	}
}

func getLogLevel(level string) zap.AtomicLevel {
	var atl zap.AtomicLevel
	switch strings.ToLower(level) {
	case DebugLevelString:
		atl = zap.NewAtomicLevelAt(zap.DebugLevel)
	case InfoLevelString:
		atl = zap.NewAtomicLevelAt(zap.InfoLevel)
	case WarnLevelString:
		atl = zap.NewAtomicLevelAt(zap.WarnLevel)
	case ErrorLevelString:
		atl = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		atl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return atl
}

func Debug(args ...interface{}) {
	fmt.Println("Debug", args, Zap)
	Zap.Sugar().Debug(args...)
}

func Debugf(format string, v ...interface{}) {
	Zap.Sugar().Debugf(format, v...)
}

func Debugln(args ...interface{}) {
	Zap.Sugar().Debugln(args...)
}

func Info(args ...interface{}) {
	Zap.Sugar().Info(args...)
}

func Infof(format string, v ...interface{}) {
	Zap.Sugar().Infof(format, v...)
}

func Warn(args ...interface{}) {
	Zap.Sugar().Warn(args...)
}

func Warnf(format string, v ...interface{}) {
	Zap.Sugar().Warnf(format, v...)
}

func Error(args ...interface{}) {
	Zap.Sugar().Error(args...)
}

func Errorf(format string, v ...interface{}) {
	Zap.Sugar().Errorf(format, v...)
}

func Fatal(args ...interface{}) {
	Zap.Sugar().Fatal(args...)
}
