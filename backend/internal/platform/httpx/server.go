package httpx

import "github.com/gin-gonic/gin"

type Server struct {
	engine  *gin.Engine
	address string
}

func NewServer(address string) *Server {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

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
