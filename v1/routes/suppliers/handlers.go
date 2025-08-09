package suppliers

import (
	"log/slog"
	"net/http"

	"github.com/MattyMcF4tty/InventoryManager-backend/v1/schemas"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/utils"
	"github.com/gin-gonic/gin"
)

// protectedFields contains fields that the user should not be able to modify
var protectedFields = []string{"id", "created_at", "updated_at", "deleted_at"}

func GetSupplierHandler(context *gin.Context) {
	id, err := utils.GetIdFromContext(context)
	if err != nil {
		slog.Error("Failed to get ID from context", "error", err)
		context.JSON(http.StatusBadRequest, schemas.ApiResponse{
			Success: false,
			Message: "Invalid ID",
		})
		return
	}

	item, err := GetSupplier(id)
	if err != nil {
		if utils.IsCustomError(err) {
			customErr := err.(*schemas.CustomError)
			slog.Error("Failed to retrieve supplier", "id", id, "error", customErr.Details)
			context.JSON(customErr.Code, schemas.ApiResponse{
				Success: false,
				Message: customErr.Message,
			})
			return
		}

		slog.Error("Unexpected error when retrieving item", "id", id, "error", err)
		context.JSON(http.StatusInternalServerError, schemas.ApiResponse{
			Success: false,
			Message: "Failed to retrieve supplier",
		})
		return
	}

	context.JSON(http.StatusOK, schemas.ApiResponse{
		Success: true,
		Message: "Supplier retrieved successfully",
		Data:    item,
	})
}
