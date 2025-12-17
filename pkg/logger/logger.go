package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogLevel string

const (
	LogLevelError LogLevel = "error"
	LogLevelWarn  LogLevel = "warn"
	LogLevelInfo  LogLevel = "info"
	LogLevelDebug LogLevel = "debug"
)

type Logger struct {
	l *zap.SugaredLogger
}

// New creates production-ready logger
func New(level LogLevel) (*Logger, error) {
	zapLevel := parseLevel(level)

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      zapcore.OmitKey,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeLevel:    levelEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	logger := zap.New(core)
	return &Logger{l: logger.Sugar()}, nil
}

func parseLevel(level LogLevel) zapcore.Level {
	switch strings.ToLower(string(level)) {
	case "error":
		return zapcore.ErrorLevel
	case "warn":
		return zapcore.WarnLevel
	case "debug":
		return zapcore.DebugLevel
	default:
		return zapcore.InfoLevel
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05"))
}

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + strings.ToUpper(l.String()) + "]")
}

func (l *Logger) Debug(args ...interface{}) {
	l.l.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.l.Debugf(format, args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.l.Debugw(msg, keysAndValues...)
}

func (l *Logger) Info(args ...interface{}) {
	l.l.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.l.Infof(format, args...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.l.Infow(msg, keysAndValues...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.l.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.l.Warnf(format, args...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.l.Warnw(msg, keysAndValues...)
}

func (l *Logger) Error(args ...interface{}) {
	l.l.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.l.Errorf(format, args...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.l.Errorw(msg, keysAndValues...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.l.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.l.Fatalf(format, args...)
}