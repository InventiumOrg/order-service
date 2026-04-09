package api

import (
	"context"
	"log/slog"
	"order-service/observability"
	routes "order-service/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Server struct {
	router            *gin.Engine
	routes            *routes.Route
	db                *pgx.Conn
	otelShutdown      func(context.Context) error
	metrics           *observability.AppMetrics
	prometheusMetrics *observability.PrometheusMetrics
}

func NewServer(db *pgx.Conn, serviceName, serviceVersion, otelEndpoint, otelHeaders string) *Server {
	server := &Server{
		router: gin.Default(),
		db:     db,
	}

	// Initialize OpenTelemetry
	ctx := context.Background()
	otelShutdown, err := observability.SetupOTelSDK(ctx, serviceName, serviceVersion, otelEndpoint, otelHeaders)
	if err != nil {
		slog.Error("Failed to setup OpenTelemetry", slog.Any("err", err))
	} else {
		server.otelShutdown = otelShutdown
		slog.Info("OpenTelemetry initialized successfully",
			slog.String("service", serviceName),
			slog.String("endpoint", otelEndpoint))
	}

	// Initialize metrics
	metrics, err := observability.CreateMetrics()
	if err != nil {
		slog.Error("Failed to create metrics", slog.Any("err", err))
	} else {
		server.metrics = metrics
		slog.Info("Metrics initialized successfully")
	}

	// Initialize Prometheus metrics
	prometheusMetrics := observability.NewPrometheusMetrics(serviceName)
	server.prometheusMetrics = prometheusMetrics
	slog.Info("Prometheus metrics initialized successfully")

	server.routes = routes.NewRoute(db, prometheusMetrics)
	return server
}

func (s *Server) Run(addr string, serviceName string) error {
	_ = s.router.SetTrustedProxies(nil)

	// Add OpenTelemetry middleware
	s.router.Use(otelgin.Middleware(serviceName))

	// Add Prometheus metrics middleware
	if s.prometheusMetrics != nil {
		s.router.Use(s.prometheusMetrics.PrometheusMiddleware())
	}

	// Add metrics middleware (OpenTelemetry)
	if s.metrics != nil {
		s.router.Use(s.metricsMiddleware())
	}

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origins", "Content-Type", "Authorization", "Bearer"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup Prometheus /metrics endpoint
	observability.SetupPrometheusEndpoint(s.router)

	// Add health check routes (no auth required)
	s.routes.AddHealthRoutes(s.router)

	// Add business logic routes
	s.routes.AddOrderRoutes(s.router)

	return s.router.Run(addr)
}

// Shutdown gracefully shuts down the server and OpenTelemetry
func (s *Server) Shutdown(ctx context.Context) error {
	if s.otelShutdown != nil {
		if err := s.otelShutdown(ctx); err != nil {
			slog.Error("Failed to shutdown OpenTelemetry", slog.Any("err", err))
			return err
		}
		slog.Info("OpenTelemetry shutdown successfully")
	}
	return nil
}

// metricsMiddleware records HTTP request metrics
func (s *Server) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()

		// Record request counter
		s.metrics.RequestCounter.Add(c.Request.Context(), 1,
			metric.WithAttributes(
				attribute.String("method", c.Request.Method),
				attribute.String("route", c.FullPath()),
				attribute.Int("status_code", c.Writer.Status()),
			))

		// Record request duration
		s.metrics.RequestDuration.Record(c.Request.Context(), duration,
			metric.WithAttributes(
				attribute.String("method", c.Request.Method),
				attribute.String("route", c.FullPath()),
				attribute.Int("status_code", c.Writer.Status()),
			))
	}
}
