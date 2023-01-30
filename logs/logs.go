package logs

import (
	slog "log"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	DefaultLogger Logger
)

type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

func SetLogger(l Logger) {
	DefaultLogger = l
}

func Debugf(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Debugf(format, v...)
	}
}

func Infof(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Infof(format, v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Warnf(format, v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Errorf(format, v...)
	}
}

func Fatalf(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Fatalf(format, v...)
	}
}
func Printf(format string, v ...interface{}) {
	if DefaultLogger != nil {
		DefaultLogger.Debugf(format, v...)
	}
}

func NewDefault(level int, envMode string) (*zap.Logger, error) {
	var lcf zap.Config
	if strings.ToLower(envMode) == "local" {
		lcf = zap.NewDevelopmentConfig()
	} else {
		lcf = zap.NewProductionConfig()
	}

	lcf.Development = true
	lcf.DisableStacktrace = true
	if strings.ToLower(envMode) == "prod" {
		lcf.Development = false
		lcf.DisableStacktrace = true
	}
	lcf.Level.SetLevel(zapcore.Level(level))
	return lcf.Build(zap.AddCallerSkip(1))
}

func InitDefaultLogger(level int, envMode string) {
	logger, err := NewDefault(level, envMode)
	if err != nil {
		slog.Fatalln("Cannot init default logger:", err.Error())
	}
	SetLogger(logger.Sugar())
}
