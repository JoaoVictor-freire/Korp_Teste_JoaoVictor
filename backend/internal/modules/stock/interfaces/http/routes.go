package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"korp_backend/internal/platform/swaggerx"
)

func RegisterRoutes(engine *gin.Engine, handler Handler, authMiddleware gin.HandlerFunc) {
	swaggerx.Register(engine, "/swagger", stockSwaggerDoc)

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "stock-service",
			"status":  "ok",
		})
	})

	v1 := engine.Group("/api/v1")
	v1.Use(authMiddleware)
	products := v1.Group("/products")
	products.POST("", handler.CreateProduct)
	products.GET("", handler.ListProducts)
	products.PUT("/:code", handler.UpdateProduct)
	products.DELETE("/:code", handler.DeleteProduct)
}
