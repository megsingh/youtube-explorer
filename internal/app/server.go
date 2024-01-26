package app

import (
	"os"
	"youtube_project/internal/api"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	videoService api.VideoService
}

func NewServer(router *gin.Engine, videoService api.VideoService) *Server {
	return &Server{
		router:       router,
		videoService: videoService,
	}
}

// Run starts the server by initializing routes and running the Gin router.
func (server *Server) Run() error {
	// run function that initializes the routes
	router := server.Routes()

	// run the server through the router
	err := router.Run(os.Getenv("PORT"))
	if err != nil {
		return err
	}

	return nil
}
