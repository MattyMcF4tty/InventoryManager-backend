package v1

import (
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/routes"
	"github.com/gin-gonic/gin"
)

func RouteHandler(v1Routes *gin.RouterGroup) {

	itemRoutes := v1Routes.Group("/items")
	routes.SetupItemRoutes(itemRoutes)
}
