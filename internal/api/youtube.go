package api

import (
	"context"
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

// RenewServiceAPIKey increments the API key index, renewing the service with the next available API key.
func (youtubeService *YouTubeAPIService) RenewServiceAPIKey() {

	apiKeys := strings.Split(os.Getenv("API_KEY"), " | ")
	youtubeService.ApiKeyId += 1
	youtubeService.ApiKey = apiKeys[youtubeService.ApiKeyId]
}

// GetLatestVideos fetches the latest videos from YouTube for a given search query.
func (youtubeService *YouTubeAPIService) GetLatestVideos(searchQuery string) ([]models.Video, error) {

	// create a new YouTube service with the current API key.
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(youtubeService.ApiKey))
	if err != nil {
		return nil, err
	}

	// set the publishedAfter parameter to the current time minus 1 year
	publishedAfter := time.Now().Add(24 * 365 * time.Hour).Format(time.RFC3339)

	// define the YouTube API call parameters.
	call := service.Search.List([]string{"snippet"}).
		Q(searchQuery).
		Type("video").
		PublishedAfter(publishedAfter).
		Order("date").
		MaxResults(10)

	response, err := call.Do()
	if err != nil {
		return nil, err
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
			return nil, err
		}
		video.PublishedAt = publishedAtTime
		video.ThumbnailUrl = item.Snippet.Thumbnails.Default.Url
		newVideos = append(newVideos, video)
	}
	return newVideos, nil
}
