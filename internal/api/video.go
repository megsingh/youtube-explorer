package api

import (
	"os"
	"youtube_project/internal/models"
	"youtube_project/internal/repository"
)

// VideoService provides methods for interacting with video data
type VideoService interface {
	InsertInDB(video []models.Video) error
	FetchFromYoutube(youtubeService YouTubeAPIService) ([]models.Video, error)
	GetAll(token string) (models.PaginationResponse, error)
	QueryDB(query, token string) (models.PaginationResponse, error)
}

// videoService implements VideoService
type videoService struct {
	storage repository.Storage
}

func NewVideoService(videoRepo repository.Storage) VideoService {
	return &videoService{
		storage: videoRepo,
	}
}

func (v *videoService) InsertInDB(videos []models.Video) error {

	err := v.storage.InsertVideos(videos)
	if err != nil {
		return err
	}

	return nil
}

func (v *videoService) FetchFromYoutube(youtubeService YouTubeAPIService) ([]models.Video, error) {

	searchQuery := os.Getenv("SEARCH_QUERY")
	videos, err := youtubeService.GetLatestVideos(searchQuery)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (v *videoService) GetAll(nextToken string) (models.PaginationResponse, error) {

	response, err := v.storage.GetPaginatedVideos(nextToken)
	if err != nil {
		return models.PaginationResponse{}, err
	}

	return response, nil
}

func (v *videoService) QueryDB(query, nextToken string) (models.PaginationResponse, error) {

	response, err := v.storage.SearchVideos(query, nextToken)
	if err != nil {
		return models.PaginationResponse{}, err
	}

	return response, nil
}
