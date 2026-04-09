package utils

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PathOrderID parses the :id route parameter. On failure it logs, writes 400, and returns ok=false.
func PathOrderID(ctx *gin.Context, rejectPrefix string) (id int64, ok bool) {
	idStr := ctx.Param("id")
	var err error
	id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Info(rejectPrefix+": invalid order id", slog.String("id_param", idStr), slog.Any("err", err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return 0, false
	}
	return id, true
}

// OrderFormInts reads pos_id, price, recipe_id from form. On failure it logs, writes 400, and returns ok=false.
// orderID is optional; when non-nil it is included in rejection logs (e.g. update).
func OrderFormInts(ctx *gin.Context, rejectPrefix string, orderID *int64) (posID, price, recipeID int32, ok bool) {
	posIDStr := ctx.PostForm("pos_id")
	priceStr := ctx.PostForm("price")
	recipeIDStr := ctx.PostForm("recipe_id")

	if posIDStr == "" || priceStr == "" || recipeIDStr == "" {
		args := []any{
			slog.Bool("has_pos_id", posIDStr != ""),
			slog.Bool("has_price", priceStr != ""),
			slog.Bool("has_recipe_id", recipeIDStr != ""),
		}
		if orderID != nil {
			args = append([]any{slog.Int64("order.id", *orderID)}, args...)
		}
		slog.Info(rejectPrefix+": missing parameters", args...)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters: pos_id, price, recipe_id",
		})
		return 0, 0, 0, false
	}

	pID, err := strconv.ParseInt(posIDStr, 10, 32)
	if err != nil {
		logInvalidFormField(ctx, rejectPrefix, "invalid pos_id", orderID, slog.String("pos_id", posIDStr), err)
		return 0, 0, 0, false
	}

	p, err := strconv.ParseInt(priceStr, 10, 32)
	if err != nil {
		logInvalidFormField(ctx, rejectPrefix, "invalid price", orderID, slog.String("price", priceStr), err)
		return 0, 0, 0, false
	}

	rID, err := strconv.ParseInt(recipeIDStr, 10, 32)
	if err != nil {
		logInvalidFormField(ctx, rejectPrefix, "invalid recipe_id", orderID, slog.String("recipe_id", recipeIDStr), err)
		return 0, 0, 0, false
	}

	return int32(pID), int32(p), int32(rID), true
}

func logInvalidFormField(ctx *gin.Context, rejectPrefix, reason string, orderID *int64, fieldAttr slog.Attr, err error) {
	args := []any{slog.String("reason", reason), fieldAttr, slog.Any("err", err)}
	if orderID != nil {
		args = append([]any{slog.Int64("order.id", *orderID)}, args...)
	}
	slog.Info(rejectPrefix+": "+reason, args...)
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": invalidFormClientError(reason),
	})
}

func invalidFormClientError(reason string) string {
	switch reason {
	case "invalid pos_id":
		return "Invalid pos_id"
	case "invalid price":
		return "Invalid price"
	case "invalid recipe_id":
		return "Invalid recipe_id"
	default:
		return "Invalid request"
	}
}
