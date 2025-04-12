package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

func main() {
	// Configure JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Basic log messages with different levels
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Log with string attributes
	logger.Info("User logged in",
		slog.String("username", "johndoe"),
		slog.String("ip", "192.168.1.1"))

	// Log with numeric attributes
	logger.Info("API request completed",
		slog.Int("status_code", 200),
		slog.Int64("response_time_ms", 127),
		slog.Float64("response_size_kb", 24.7))

	// Log with boolean attributes
	logger.Info("Feature flags status",
		slog.Bool("dark_mode", true),
		slog.Bool("beta_features", false),
		slog.Any("thingy", nil))

	fmt.Println("something without proper JSON")

	// Log with complex JSON structure using Any
	user := map[string]any{
		"id":       12345,
		"username": "alice",
		"roles":    []string{"admin", "user"},
		"settings": map[string]any{
			"notifications": true,
			"theme":         "dark",
			"thing":         nil,
		},
	}
	logger.Info("User profile", slog.Any("user", user))

	// Log with array/slice
	tags := []string{"go", "logging", "slog", "json"}
	logger.Info("Article published",
		slog.String("title", "Structured Logging in Go"),
		slog.Any("tags", tags))

	// Log with time
	logger.Info("Event scheduled",
		slog.Time("scheduled_at", time.Now().Add(24*time.Hour)))

	// Log with nested error
	logger.Error("Database connection failed",
		slog.String("db", "users"),
		slog.Any("error", map[string]any{
			"code":    "CONN_REFUSED",
			"message": "connection refused",
			"details": map[string]any{
				"host":        "db.example.com",
				"port":        5432,
				"retry_after": 30,
				"foo":         true,
				"bar":         false,
			},
		}))

	// Mixed types in a single log
	logger.Info("Order processed",
		slog.Int("order_id", 987654),
		slog.String("customer", "ACME Corp"),
		slog.Float64("total", 299.99),
		slog.Bool("expedited", true),
		slog.Any("items", []map[string]any{
			{
				"product_id": "ABC123",
				"quantity":   2,
				"price":      149.99,
			},
			{
				"product_id": "XYZ789",
				"quantity":   1,
				"price":      0.01,
			},
		}))

	badLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.MessageKey:
				a.Key = "message"
			case slog.TimeKey:
				a.Key = "timestamp"
			case slog.LevelKey:
				a.Key = "logLevel"
			}
			return a
		},
	}))

	badLogger.Info("big chaos",
		slog.Float64("total", 299.99),
		slog.Bool("expedited", true),
		slog.Any("items", []map[string]any{
			{
				"product_id": "ABC123",
				"quantity":   2,
				"price":      149.99,
			},
			{
				"product_id": "XYZ789",
				"quantity":   1,
				"price":      0.01,
			},
		}))
}
