package api

import (
	"AutoScaler/internal/api/handler"
	"AutoScaler/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	storage storage.Storage
	router  *gin.Engine
}

func NewServer(storage storage.Storage) *Server {
	router := gin.Default()

	s := &Server{
		storage: storage,
		router:  router,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {

	groupHandler := handler.NewGroupHandler(s.storage)

	s.router.GET("/healthz", s.healthCheck)

	groups := s.router.Group("/v1/groups")
	{
		groups.POST("", groupHandler.CreateGroup)
		groups.GET("", groupHandler.ListGroups)
		groups.GET("/:name", groupHandler.GetGroup)
		groups.PUT("/:name", groupHandler.UpdateGroup)
		groups.DELETE("/:name", groupHandler.DeleteGroup)
		groups.PUT("/:name/desired-capacity", groupHandler.SetDesiredCapacity)
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
