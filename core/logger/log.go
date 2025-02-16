package logger

import (
	"fmt"
	"github.com/darkit/slog"
	"os"
)

func InitLogger(_ string) {
	slog.NewLogger(os.Stdout, true, false)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Error(args ...any) {
	slog.Errorf("%s", fmt.Sprintf("%v", args...))
}

func Trace(msg string, args ...any) {
	slog.Trace(msg, args...)
}

func Tracef(format string, args ...any) {
	slog.Tracef(format, args...)
}

func Debugf(format string, args ...any) {
	slog.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	slog.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	slog.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	slog.Errorf(format, args...)
}

func Fatal(args ...any) {
	slog.Fatalf("%s", fmt.Sprintf("%v", args...))
}

func Fatalf(format string, args ...any) {
	slog.Fatalf(format, args...)
}
