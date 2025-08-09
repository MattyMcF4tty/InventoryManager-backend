package suppliers

import (
	"github.com/gin-gonic/gin"
)

func SetupSupplierRoutes(routes *gin.RouterGroup) {
	routes.GET("/:id", GetSupplierHandler)
}
