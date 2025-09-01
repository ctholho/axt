package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	// Create a new production-ready Zap logger.
	// This logger outputs logs in JSON format.
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync() // Flushes buffer, if any.

	// Use a sugared logger for easier, more conversational logging
	// if needed. We will stick with the standard logger for structured logs.
	// sugar := logger.Sugar()

	// 1. Simulate server startup
	logger.Info("Server application starting.",
		zap.String("service", "web-server"),
		zap.String("version", "1.0.0"),
		zap.String("environment", "development"),
	)

	// 2. Simulate an incoming request
	requestId := "d3a2b4f" // A simulated request ID
	logger.Info("Incoming request received.",
		zap.String("method", "GET"),
		zap.String("path", "/api/v1/users"),
		zap.String("client_ip", "127.0.0.1"),
		zap.String("request_id", requestId),
	)

	// 3. Simulate a warning event
	logger.Warn("Database connection is slow.",
		zap.Int("duration_ms", 250),
		zap.String("database", "user_db"),
	)

	// 4. Simulate a business logic event
	logger.Info("User registration successful.",
		zap.String("user_id", "user-1234"),
		zap.String("source", "web-form"),
	)

	// 5. Simulate an error event
	logger.Error("Failed to write to file.",
		zap.String("file_path", "/var/log/app.log"),
		zap.String("error", "Permission denied"),
		zap.String("user", "app-user"),
	)

	// 6. Simulate server shutdown
	logger.Info("Server gracefully shutting down.",
		zap.String("reason", "idle_timeout"),
	)

	// The program exits here, and the defer will run, flushing the logs.
}
