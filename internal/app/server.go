package app

import (
	"log"
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

func (server *Server) Run() error {
	// run function that initializes the routes
	r := server.Routes()

	// run the server through the router
	err := r.Run(os.Getenv("PORT"))

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}
