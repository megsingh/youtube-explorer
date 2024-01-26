package app

import (
	"net/http"
	"youtube_project/pkg/error_handler/api_errors"

	"github.com/gin-gonic/gin"
)

// ApiStatus is the api for "/" route
func (s *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "youtube video search API running smoothly"})
	}
}

// GetAllVideos returns a Gin handler function that fetches all videos.
func (server *Server) GetAllVideos() gin.HandlerFunc {
	return func(c *gin.Context) {

		nextToken := c.Query("next_token")

		// Call the videoService to fetch all videos with pagination using the provided 'nextToken'.
		responseVideos, err := server.videoService.GetAll(nextToken)
		if err != nil {
			fetchError := api_errors.NewVideoFetchError(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fetchError.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": responseVideos})
	}
}

// Search returns a Gin handler function that searches videos based on a search query.
func (server *Server) Search() gin.HandlerFunc {
	return func(c *gin.Context) {

		nextToken := c.Query("next_token")
		searchQuery := c.Query("search_query")

		// Call the videoService to search videos in the database based on the provided 'searchQuery' and 'nextToken'.
		responseVideos, err := server.videoService.QueryDB(searchQuery, nextToken)
		if err != nil {
			fetchError := api_errors.NewVideoFetchError(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fetchError.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": responseVideos})
	}
}
