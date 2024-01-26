package app

import "github.com/gin-gonic/gin"

func (server *Server) Routes() *gin.Engine {
	router := server.router

	// group all routes under /v1/api
	v1 := router.Group("/v1/api")
	{
		v1.GET("/", server.ApiStatus())
		v1.GET("/videos", server.GetAllVideos())
		v1.GET("/search", server.Search())

	}

	return router
}
