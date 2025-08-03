package main

import (
	v1 "github.com/MattyMcF4tty/InventoryManager-backend/v1"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Get the api v1 routes
	v1Routes := router.Group("/v1")
	v1.RouteHandler(v1Routes)

	// Start server on port 8080
	router.Run("0.0.0.0:8080")
}
