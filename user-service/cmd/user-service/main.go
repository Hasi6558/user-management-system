package main

import (
	"log/slog"
	"os"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler).With("service", "user-service")
	slog.SetDefault(logger)
	logger.Info("starting user service")
}
