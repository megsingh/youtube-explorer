package app

import "github.com/gin-gonic/gin"

// Routes configures and returns the Gin router with all the API endpoints.
func (server *Server) Routes() *gin.Engine {

	// group all routes under /v1/api
	v1 := server.router.Group("/v1/api")
	{
		v1.GET("/", server.ApiStatus())
		v1.GET("/videos", server.GetAllVideos())
		v1.GET("/search", server.Search())

	}

	return server.router
}
