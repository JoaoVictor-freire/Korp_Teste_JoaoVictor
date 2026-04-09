package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine, handler Handler) {
	v1 := engine.Group("/api/v1")

	authGroup := v1.Group("/auth")
	authGroup.POST("/register", handler.Register)
	authGroup.POST("/login", handler.Login)

	// Convenience endpoint for demos.
	engine.GET("/whoami", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hint": "use /api/v1/auth/login and Authorization: Bearer <token> on /api/v1/* endpoints",
		})
	})
}
