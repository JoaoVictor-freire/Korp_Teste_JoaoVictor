package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"korp_backend/internal/platform/swaggerx"
)

func RegisterRoutes(engine *gin.Engine, handler Handler, authMiddleware gin.HandlerFunc) {
	swaggerx.Register(engine, "/swagger", billingSwaggerDoc)

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
	invoices.PATCH("/:number/close", handler.CloseInvoice)
}
