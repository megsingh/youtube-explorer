package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) ApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"data": "youtube video search API running smoothly",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) GetAllVideos() gin.HandlerFunc {
	return func(c *gin.Context) {

		nextToken := c.Query("next_token")
		videos, err := s.videoService.GetAll(nextToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch videos"})
			return
		}

		response := map[string]interface{}{
			"data": videos,
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) Search() gin.HandlerFunc {
	return func(c *gin.Context) {

		nextToken := c.Query("next_token")
		search_query := c.Query("search_query")

		videos, err := s.videoService.SearchAll(search_query, nextToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search videos"})
			return
		}

		response := map[string]interface{}{
			"data": videos,
		}
		c.JSON(http.StatusOK, response)
	}
}
