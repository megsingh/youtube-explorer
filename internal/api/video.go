package api

import (
	"log"
	"os"
	"youtube_project/internal/models"
	"youtube_project/internal/repository"
)

// VideoService provides methods for interacting with video data
type VideoService interface {
	Insert(video []models.Video) error
	FetchFromYoutube(youtubeService YouTubeAPIService) ([]models.Video, error)
	GetAll(token string) (models.PaginationResponse, error)
	SearchAll(query, token string) (models.PaginationResponse, error)
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

func (v *videoService) Insert(videos []models.Video) error {
	log.Println("inserting videos")
	err := v.storage.InsertVideos(videos)

	if err != nil {
		return err
	}

	log.Println("videos inserted successfully")
	return nil
}

func (v *videoService) FetchFromYoutube(youtubeService YouTubeAPIService) ([]models.Video, error) {
	log.Println("fetching videos")
	searchQuery := os.Getenv("SEARCH_QUERY")
	videos, err := youtubeService.GetLatestVideos(searchQuery)
	if err != nil {
		return []models.Video{}, err
	}

	log.Println("videos fetched from youtube successfully")

	return videos, nil
}

func (v *videoService) GetAll(searchToken string) (models.PaginationResponse, error) {
	log.Println("getting all videos according to published date")

	response, err := v.storage.GetPaginatedVideos(searchToken)
	if err != nil {
		return models.PaginationResponse{}, err
	}

	log.Println("all videos fetched from dbsuccessfully")

	return response, nil
}

func (v *videoService) SearchAll(query, nextToken string) (models.PaginationResponse, error) {
	log.Println("searching all videos according to query")

	response, err := v.storage.SearchVideos(query, nextToken)
	if err != nil {
		return models.PaginationResponse{}, err
	}

	log.Println("search successful")

	return response, nil
}
