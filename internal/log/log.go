package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// global variable Logger to avoid log injections everywhere
var (
	esLogger     *zap.Logger
	levelChanger zap.AtomicLevel
)

func init() {
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)
	l := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout),
		atom,
	))
	SetupLogger(l, atom)
}

func SetupLogger(l *zap.Logger, _levelChanger zap.AtomicLevel) {
	esLogger = l.WithOptions(zap.AddCallerSkip(1))
	levelChanger = _levelChanger
}

func GetLogger() *zap.Logger {
	return esLogger
}

func SetLogLevel(level zapcore.Level) {
	levelChanger.SetLevel(level)
}

func GetLogLevel() zapcore.Level {
	return levelChanger.Level()
}

func Named(s string) *zap.Logger {
	return esLogger.Named(s)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return esLogger.WithOptions(opts...)
}

func With(fields ...zap.Field) *zap.Logger {
	return esLogger.With(fields...)
}

func Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return esLogger.Check(lvl, msg)
}

func Debug(msg string, fields ...zap.Field) {
	esLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	esLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	esLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	esLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	esLogger.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	esLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	esLogger.Fatal(msg, fields...)
}

func Sync() error {
	return esLogger.Sync()
}

func Core() zapcore.Core {
	return esLogger.Core()
}
