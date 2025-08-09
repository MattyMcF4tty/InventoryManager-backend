package items

import (
	"github.com/gin-gonic/gin"
)

func SetupItemRoutes(routes *gin.RouterGroup) {
	routes.GET("", GetPagedItemsHandler)
	routes.GET("/:id", GetItemHandler)
	routes.PATCH("/:id", UpdateItemHandler)
	routes.POST("/", CreateItemHandler)
	routes.DELETE("/:id", DeleteItemHandler)
}
