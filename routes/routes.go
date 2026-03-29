package routes

import (
	handlers "order-service/handlers"
	"order-service/observability"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Route struct {
	db       *pgx.Conn
	handlers *handlers.Handlers
}

func NewRoute(db *pgx.Conn, prometheusMetrics *observability.PrometheusMetrics) *Route {
	return &Route{
		db:       db,
		handlers: handlers.NewHandlers(db, prometheusMetrics),
	}
}

func (r *Route) AddOrderRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.GET("/:id", r.handlers.GetOrder)
			orders.GET("/list", r.handlers.ListOrder)
			orders.POST("/create", r.handlers.CreateOrder)
			orders.PUT("/:id", r.handlers.UpdateOrder)
			orders.DELETE("/:id", r.handlers.DeleteOrder)
		}
	}
}

func (r *Route) AddHealthRoutes(router *gin.Engine) {
	// Health check endpoints (no authentication required)
	router.GET("/healthz", r.handlers.HealthzHandler)
	router.GET("/readyz", r.handlers.ReadyzHandler)
}
