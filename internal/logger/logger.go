package logger

import (
	"log/slog"
	"os"
)

func CreateNewLogger(serviceName string) *slog.Logger {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	l = l.With("Service Name", serviceName)
	return l
}
