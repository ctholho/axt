//go:build debug

package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Basic log messages with different levels
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Log strings properties
	logger.Info("User logged in",
		slog.String("username", "johndoe"),
		slog.String("ip", "192.168.1.1"))

	// Log numbers and bools
	logger.Info("API request completed",
		slog.Int("status_code", 200),
		slog.Int64("response_time_ms", 127),
		slog.Float64("response_size_kb", 24.7),
		slog.Bool("beta_features", false))

	fmt.Println("something without proper JSON")

	// Log complex JSON structure
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

	// Log array
	tags := []string{"go", "logging", "slog", "json"}
	logger.Info("Article published",
		slog.String("title", "Structured Logging in Go"),
		slog.Any("tags", tags))

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

	// Extremely long line
	logger.Info("Order processed but nothing good ever comes from too many orders. All in all one could say that Marx was right when he said: Do I obey economic laws if I extract money by offering my body for sale,... — Then the political economist replies to me: You do not transgress my laws; but see what Cousin Ethics and Cousin Religion have to say about it. My political economic ethics and religion have nothing to reproach you with, but — But whom am I now to believe, political economy or ethics? — The ethics of political economy is acquisition, work, thrift, sobriety — but political economy promises to satisfy my needs. ... It stems from the very nature of estrangement that each sphere applies to me a different and opposite yardstick — ethics one and political economy another; for each is a specific estrangement of man and focuses attention on a particular field of estranged essential activity, and each stands in an estranged relation to the other.",
		slog.String("customer", "The International Consortium for the Global Provision of Highly Specialized, Contextually-Appropriate, and Sustainably-Sourced Digital and Analog Solutions for the Modern Enterprise, Including but Not Limited to Cloud-Based Infrastructure, Bespoke Software Development, Strategic Marketing Campaigns, and High-Fidelity Audio-Visual Production Services, a Subsidiary of the Consolidated Federation of Intercontinental Business Ventures and Allied Partners, LLC."),
		slog.Float64("total", 0.99),
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

	// Epoch time to test --time-in Unix (use with -t timestamp)
	logger.Info("Epoch Time",
		slog.String("timestamp", "1756555555"),
	)
}
