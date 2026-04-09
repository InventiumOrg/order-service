package main

import (
	"context"
	"log/slog"
	"order-service/api"
	"order-service/config"
	"order-service/observability"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

const attemptThreshold = 5

// setupLogging configures logging based on environment variables
func setupLogging(cfg config.Config) error {

	// Use OTLP for logs if configured
	if cfg.OTELExporterOTLPEndpoint != "" {
		endpoint := "http://" + cfg.OTELExporterOTLPEndpoint
		if err := observability.SetupOTLPLogging(endpoint, cfg.ServiceName); err == nil {
			slog.Info("OTLP logging configured", slog.String("endpoint", endpoint))
			return nil
		} else {
			slog.Warn("OTLP logging failed, falling back to stdout", slog.Any("error", err))
		}
	}

	slog.Info("Using default stdout logging")
	return nil
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("Failed to load config: ", slog.Any("ERROR", err))
		os.Exit(1)
	}
	slog.Info("Set Up Logging.....")

	// Setup logging based on configuration
	if err := setupLogging(config); err != nil {
		slog.Error("Failed to setup logging", slog.Any("error", err))
		// Continue with stdout logging if setup fails
	}
	time.Sleep(10 * time.Second)
	slog.Info("Connecting to database...")
	attempt := 1
	for attempt <= attemptThreshold {
		conn, err = pgx.Connect(context.Background(), config.DBSource)
		if err == nil {
			slog.Info("Connected to database successfully")
			// defer conn.Close(context.Background())
			break
		}
		slog.Error("Failed to connect to database",
			slog.Int("attempt", attempt),
			slog.Int("maxAttempts", attemptThreshold),
			slog.Any("error", err),
		)

		if attempt == attemptThreshold {
			slog.Error("Max connection attempts reached, exiting", slog.Any("ERROR", err))
			os.Exit(1)
		}

		backoffDuration := time.Duration(1<<(attempt-1)) * time.Second
		slog.Info("Retrying connection",
			slog.Int("attempt", attempt+1),
			slog.Duration("backoff", backoffDuration),
		)

		time.Sleep(backoffDuration)
		attempt++

	}
	router := api.NewServer(conn, config.ServiceName, "1.0.0", config.OTELExporterOTLPEndpoint, config.OTELExporterOTLPHeaders)
	_ = router.Run(":9820", config.ServiceName)

}
