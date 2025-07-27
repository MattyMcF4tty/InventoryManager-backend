package routes

import (
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/items"
	"github.com/gin-gonic/gin"
)

func SetupItemRoutes(routes *gin.RouterGroup) {
	routes.GET("/:id", items.GetItemHandler)
	routes.PATCH("/:id", items.UpdateItemHandler)
	routes.POST("/", items.CreateItemHandler)
	routes.DELETE("/:id", items.DeleteItemHandler)
}
