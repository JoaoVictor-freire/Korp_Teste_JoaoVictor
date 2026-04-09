package swaggerx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Register(engine *gin.Engine, basePath string, docJSON string) {
	jsonPath := basePath + ".json"

	engine.GET(jsonPath, func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(docJSON))
	})

	engine.GET(basePath+"/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL(jsonPath),
	)))
}
