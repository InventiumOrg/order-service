package handlers

import (
	"log/slog"
	"net/http"
	models "order-service/models/sqlc"
	"order-service/observability"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Handlers struct {
	queries           *models.Queries
	tracer            trace.Tracer
	db                *pgx.Conn
	prometheusMetrics *observability.PrometheusMetrics
}

func NewHandlers(db *pgx.Conn, prometheusMetrics *observability.PrometheusMetrics) *Handlers {
	return &Handlers{
		db:                db,
		queries:           models.New(db),
		tracer:            otel.Tracer("order-service/handlers"),
		prometheusMetrics: prometheusMetrics,
	}
}

func (h *Handlers) GetOrder(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "GetOrder")
	defer span.End()

	// Get order ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	// Add attributes to the span
	span.SetAttributes(attribute.Int64("order.id", id))

	dbStart := time.Now()
	order, err := h.queries.GetOrder(ctx, id)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("get", "orders", dbDuration, err)
	}

	if err != nil {
		slog.Error("Got an error while getting order", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get order",
		})
		return
	}

	// Record successful retrieval (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordOrderOperation("get", order.PosID)
	}

	// Record successful operation
	span.SetAttributes(
		attribute.Int("order.pos_id", int(order.PosID)),
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(200, gin.H{
		"message": "Get Order Successfully",
		"data":    order,
	})
}

func (h *Handlers) ListOrder(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "ListOrders")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int("order.limit", 10),
		attribute.Int("order.offset", 0),
	)

	dbStart := time.Now()
	orders, err := h.queries.ListOrder(ctx, models.ListOrderParams{
		Limit:  10,
		Offset: 0,
	})
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("list", "orders", dbDuration, err)
	}

	if err != nil {
		slog.Error("Got an error while listing orders", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list orders",
		})
		return
	}

	// Record successful list operation (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordOrderOperation("list", 0)
	}

	span.SetAttributes(
		attribute.Int("order.count", len(orders)),
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(200, gin.H{
		"message": "List Orders Successfully",
		"data":    orders,
		"count":   len(orders),
	})
}

func (h *Handlers) CreateOrder(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "CreateOrder")
	defer span.End()

	// Parse form values
	posIDStr := ctx.PostForm("pos_id")
	priceStr := ctx.PostForm("price")
	recipeIDStr := ctx.PostForm("recipe_id")

	if posIDStr == "" || priceStr == "" || recipeIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters: pos_id, price, recipe_id",
		})
		return
	}

	posID, err := strconv.ParseInt(posIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid pos_id",
		})
		return
	}

	price, err := strconv.ParseInt(priceStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid price",
		})
		return
	}

	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid recipe_id",
		})
		return
	}

	param := models.CreateOrderParams{
		PosID:    int32(posID),
		Price:    int32(price),
		RecipeID: int32(recipeID),
	}

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int("order.pos_id", int(posID)),
		attribute.Int("order.price", int(price)),
		attribute.Int("order.recipe_id", int(recipeID)),
	)

	dbStart := time.Now()
	order, err := h.queries.CreateOrder(ctx, param)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("create", "orders", dbDuration, err)
	}

	if err != nil {
		slog.Error("Could not create order", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create order",
		})
		return
	}

	// Record successful creation (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordOrderOperation("create", order.PosID)
		h.prometheusMetrics.UpdateOrdersCount(1)
	}

	// Record successful operation
	span.SetAttributes(
		attribute.Int64("order.id", order.ID),
		attribute.String("operation.status", "success"),
	)

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Create Order Successfully",
		"data":    order,
	})
}

func (h *Handlers) UpdateOrder(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "UpdateOrder")
	defer span.End()

	// Get order ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	// Parse form values
	posIDStr := ctx.PostForm("pos_id")
	priceStr := ctx.PostForm("price")
	recipeIDStr := ctx.PostForm("recipe_id")

	if posIDStr == "" || priceStr == "" || recipeIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters: pos_id, price, recipe_id",
		})
		return
	}

	posID, err := strconv.ParseInt(posIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid pos_id",
		})
		return
	}

	price, err := strconv.ParseInt(priceStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid price",
		})
		return
	}

	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid recipe_id",
		})
		return
	}

	param := models.UpdateOrderParams{
		ID:       id,
		PosID:    int32(posID),
		Price:    int32(price),
		RecipeID: int32(recipeID),
	}

	// Add attributes to the span
	span.SetAttributes(
		attribute.Int64("order.id", id),
		attribute.Int("order.pos_id", int(posID)),
	)

	dbStart := time.Now()
	order, err := h.queries.UpdateOrder(ctx, param)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("update", "orders", dbDuration, err)
	}

	if err != nil {
		slog.Error("Could not update order", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update order",
		})
		return
	}

	// Record successful update (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordOrderOperation("update", order.PosID)
	}

	span.SetAttributes(attribute.String("operation.status", "success"))

	ctx.JSON(200, gin.H{
		"message": "Update Order Successfully",
		"data":    order,
	})
}

func (h *Handlers) DeleteOrder(ctx *gin.Context) {
	// Start a new span for this operation
	_, span := h.tracer.Start(ctx.Request.Context(), "DeleteOrder")
	defer span.End()

	// Get order ID from URL parameter
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	// Add attributes to the span
	span.SetAttributes(attribute.Int64("order.id", id))

	dbStart := time.Now()
	err = h.queries.DeleteOrder(ctx, id)
	dbDuration := time.Since(dbStart)

	// Record database operation duration (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordDBOperation("delete", "orders", dbDuration, err)
	}

	if err != nil {
		slog.Error("Could not delete order", slog.Any("err", err.Error()))
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete order",
		})
		return
	}

	// Record successful deletion (Prometheus)
	if h.prometheusMetrics != nil {
		h.prometheusMetrics.RecordOrderOperation("delete", 0)
		h.prometheusMetrics.UpdateOrdersCount(-1)
	}

	span.SetAttributes(attribute.String("operation.status", "success"))

	ctx.JSON(200, gin.H{
		"message": "Delete Order Successfully",
	})
}
