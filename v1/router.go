package v1

import (
	items "github.com/MattyMcF4tty/InventoryManager-backend/v1/routes/items"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/routes/suppliers"
	"github.com/gin-gonic/gin"
)

func RouteHandler(v1Routes *gin.RouterGroup) {

	itemRoutes := v1Routes.Group("/items")
	items.SetupItemRoutes(itemRoutes)

	supplierRoutes := v1Routes.Group("/suppliers")
	suppliers.SetupSupplierRoutes(supplierRoutes)
}
