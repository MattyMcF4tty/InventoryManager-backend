package items

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/MattyMcF4tty/InventoryManager-backend/v1/schemas"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/utils"
	"github.com/gin-gonic/gin"
)

// protectedFields contains fields that the user should not be able to modify
var protectedFields = []string{"id", "created_at", "updated_at", "deleted_at"}

func GetItemHandler(context *gin.Context) {
	id, err := utils.GetIdFromContext(context)
	if err != nil {
		slog.Error("Failed to get ID from context", "error", err)
		context.JSON(http.StatusBadRequest, schemas.ApiResponse{
			Success: false,
			Message: "Invalid ID",
		})
		return
	}

	item, err := GetItem(id)
	if err != nil {
		if utils.IsCustomError(err) {
			customErr := err.(*schemas.CustomError)
			slog.Error("Failed to retrieve item", "id", id, "error", customErr.Details)
			context.JSON(customErr.Code, schemas.ApiResponse{
				Success: false,
				Message: customErr.Message,
			})
			return
		}

		slog.Error("Unexpected error when retrieving item", "id", id, "error", err)
		context.JSON(http.StatusInternalServerError, schemas.ApiResponse{
			Success: false,
			Message: "Failed to retrieve item",
		})
		return
	}

	context.JSON(http.StatusOK, schemas.ApiResponse{
		Success: true,
		Message: "Item retrieved successfully",
		Data:    item,
	})
}

func UpdateItemHandler(context *gin.Context) {
	id, err := utils.GetIdFromContext(context)
	if err != nil {
		slog.Error("Failed to get ID from context", "error", err)
		context.JSON(http.StatusBadRequest, schemas.ApiResponse{
			Success: false,
			Message: "Invalid ID",
		})
		return
	}

	var updates map[string]interface{}
	if err := context.ShouldBindJSON(&updates); err != nil {
		slog.Error("Failed to parse JSON of updated item", "error", err)
		context.JSON(http.StatusBadRequest, schemas.ApiResponse{
			Success: false,
			Message: "Invalid JSON in body.",
		})
		return
	}

	utils.RemoveProtectedFields(updates, protectedFields)

	item, err := UpdateItem(id, updates)
	if err != nil {
		if utils.IsCustomError(err) {
			customErr := err.(*schemas.CustomError)
			slog.Error("Failed to update item", "id", id, "error", customErr.Details)
			context.JSON(customErr.Code, schemas.ApiResponse{
				Success: false,
				Message: customErr.Message,
			})
			return
		}

		slog.Error("Failed to update item", "id", id, "error", err)
		context.JSON(http.StatusInternalServerError, schemas.ApiResponse{
			Success: false,
			Message: "Failed to update item",
		})
		return
	}

	context.JSON(http.StatusOK, schemas.ApiResponse{
		Success: true,
		Message: "Item updated successfully",
		Data:    item,
	})
}

func CreateItemHandler(context *gin.Context) {
	var itemData map[string]interface{}
	if err := context.ShouldBindJSON(&itemData); err != nil {
		slog.Error("Failed to parse JSON of new item", "error", err)
		context.JSON(http.StatusBadRequest, schemas.ApiResponse{
			Success: false,
			Message: "Invalid JSON in body.",
		})
		return
	}

	err := utils.CheckRequiredFields(itemData, []string{"name", "description", "quantity", "price", "supplier_id", "category"})
	if err != nil {
		slog.Error("Missing required fields in item data", "error", err)
		context.JSON(http.StatusBadRequest, schemas.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	newItem := schemas.Item{
		Name:          itemData["name"].(string),
		Description:   itemData["description"].(string),
		Quantity:      int8(itemData["quantity"].(float64)),
		PurchasePrice: itemData["purchase_price"].(float64),
		SupplierId:    int8(itemData["supplier_id"].(float64)),
		Category:      itemData["category"].(string),
	}

	item, err := CreateItem(newItem)
	if err != nil {
		if utils.IsCustomError(err) {
			customErr := err.(*schemas.CustomError)
			slog.Error("Failed to create item", "error", customErr.Details)
			context.JSON(customErr.Code, schemas.ApiResponse{
				Success: false,
				Message: customErr.Message,
			})
			return
		}

		slog.Error("Failed to create item", "error", err)
		context.JSON(http.StatusInternalServerError, schemas.ApiResponse{
			Success: false,
			Message: "Failed to create item",
		})
		return
	}

	context.JSON(http.StatusCreated, schemas.ApiResponse{
		Success: true,
		Message: "Item created successfully",
		Data:    item,
	})
}

func DeleteItemHandler(context *gin.Context) {
	id, err := utils.GetIdFromContext(context)
	if err != nil {
		slog.Error("Failed to get ID from context", "error", err)
		context.JSON(http.StatusBadRequest, schemas.ApiResponse{
			Success: false,
			Message: "Invalid ID",
		})
		return
	}

	err = DeleteItem(id)
	if err != nil {
		if utils.IsCustomError(err) {
			customErr := err.(*schemas.CustomError)
			slog.Error("Failed to delete item", "id", id, "error", customErr.Details)
			context.JSON(customErr.Code, schemas.ApiResponse{
				Success: false,
				Message: customErr.Message,
			})
			return
		}

		slog.Error("Failed to delete item", "id", id, "error", err)
		context.JSON(http.StatusInternalServerError, schemas.ApiResponse{
			Success: false,
			Message: "Failed to delete item",
		})
		return
	}

	context.JSON(http.StatusOK, schemas.ApiResponse{
		Success: true,
		Message: "Item deleted successfully",
	})
}

func GetPagedItemsHandler(context *gin.Context) {
	pageStr := context.Query("page")
	pageSizeStr := context.Query("page-size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		context.JSON(http.StatusBadRequest, schemas.ApiResponse{
			Success: false,
			Message: "Invalid page number",
		})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		context.JSON(http.StatusBadRequest, schemas.ApiResponse{
			Success: false,
			Message: "Invalid page size",
		})
		return
	}

	items, count, err := getPagedItems(page, pageSize)
	if err != nil {
		if utils.IsCustomError(err) {
			customErr := err.(*schemas.CustomError)
			slog.Error("Failed to retrieve paged items", "error", customErr.Details)
			context.JSON(customErr.Code, schemas.ApiResponse{
				Success: false,
				Message: customErr.Message,
			})
			return
		}

		slog.Error("Failed to retrieve paged items", "error", err)
		context.JSON(http.StatusInternalServerError, schemas.ApiResponse{
			Success: false,
			Message: "Failed to retrieve paged items",
		})
		return
	}

	context.JSON(http.StatusOK, schemas.ApiResponse{
		Success: true,
		Message: "Paged items retrieved successfully",
		Data: map[string]interface{}{
			"count":    count,
			"page":     page,
			"pageSize": pageSize,
			"data":     items,
		},
	})
}
