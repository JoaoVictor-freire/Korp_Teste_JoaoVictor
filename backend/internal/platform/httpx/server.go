package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine  *gin.Engine
	address string
}

func NewServer(address string) *Server {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery(), CORSMiddleware())
	engine.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	return &Server{
		engine:  engine,
		address: address,
	}
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) Run() error {
	return s.engine.Run(s.address)
}
