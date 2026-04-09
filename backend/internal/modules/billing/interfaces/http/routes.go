package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine, handler Handler, authMiddleware gin.HandlerFunc) {
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "billing-service",
			"status":  "ok",
		})
	})

	v1 := engine.Group("/api/v1")
	v1.Use(authMiddleware)
	invoices := v1.Group("/invoices")
	invoices.POST("", handler.CreateInvoice)
	invoices.GET("", handler.ListInvoices)
}
