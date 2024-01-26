package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"youtube_project/internal/models"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeAPIService struct {
	ApiKey   string
	ApiKeyId int
}

// NewYouTubeAPIService creates a new instance of YouTubeAPIService.
func NewYouTubeAPIService() YouTubeAPIService {

	apiKeys := strings.Split(os.Getenv("API_KEY"), " | ")
	apiKeyId := 0
	return YouTubeAPIService{
		ApiKey:   apiKeys[apiKeyId],
		ApiKeyId: apiKeyId,
	}
}

func (YTService *YouTubeAPIService) RenewServiceAPIKey() {
	apiKeys := strings.Split(os.Getenv("API_KEY"), " | ")
	YTService.ApiKeyId += 1
	YTService.ApiKey = apiKeys[YTService.ApiKeyId]
}

// GetLatestVideos fetches the latest videos from YouTube for a given search query.
func (youtubeService *YouTubeAPIService) GetLatestVideos(searchQuery string) ([]models.Video, error) {
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(youtubeService.ApiKey))
	if err != nil {
		log.Println("Error creating YouTube service:", err)
		return []models.Video{}, err
	}

	// Set the publishedAfter parameter to the current time minus 24 hours
	publishedAfter := time.Now().Add(24 * 365 * time.Hour).Format(time.RFC3339)

	call := service.Search.List([]string{"snippet"}).
		Q(searchQuery).
		Type("video").
		PublishedAfter(publishedAfter).
		Order("date").
		MaxResults(10) // Adjust the number of results as needed

	response, err := call.Do()
	if err != nil {
		return []models.Video{}, err
	}

	newVideos := make([]models.Video, 0)
	for _, item := range response.Items {
		var video models.Video
		video.ChannelId = item.Snippet.ChannelId
		video.ChannelTitle = item.Snippet.ChannelTitle
		video.Title = item.Snippet.Title
		video.Description = item.Snippet.Description

		publishedAtTime, err := time.Parse("2006-01-02T15:04:05Z", item.Snippet.PublishedAt)

		if err != nil {
			fmt.Println("Error parsing date:", err)
			return []models.Video{}, err
		}
		video.PublishedAt = publishedAtTime
		video.ThumbnailUrl = item.Snippet.Thumbnails.Default.Url
		newVideos = append(newVideos, video)
	}
	return newVideos, nil
}
