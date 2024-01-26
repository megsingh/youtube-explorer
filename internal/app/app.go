package app

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"youtube_project/internal/api"
	"youtube_project/internal/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Run initializes and starts the application
func Run() error {

	// setup database connection
	dbClient, collection, err := repository.SetupDatabase()

	if err != nil {
		return err
	}

	// create storage dependency
	storage := repository.NewStorage(dbClient, collection)

	// create router dependecy
	router := gin.Default()
	router.Use(cors.Default())

	youtubeService := api.NewYouTubeAPIService()

	// create video service
	videoService := api.NewVideoService(storage)

	go func(videoService api.VideoService, youtubeService api.YouTubeAPIService) {
		for {
			videos, err := videoService.FetchFromYoutube(youtubeService)
			retryInterval, _ := strconv.Atoi(os.Getenv("API_RETRY_INTERVAL"))
			sleepInterval, _ := strconv.Atoi(os.Getenv("API_SLEEP_INTERVAL"))

			if err != nil {
				if strings.Contains(err.Error(), "quotaExceeded") {
					log.Printf("API KEY %v QUOTA EXCEEDED. RETRYING AFTER %v seconds\n", youtubeService.ApiKey, retryInterval)
					youtubeService.RenewServiceAPIKey()
					time.Sleep(time.Duration(retryInterval) * time.Second) // Retry in case of an error
					continue
				} else {
					log.Fatalln(err)
				}
			}
			err = videoService.Insert(videos)

			if err != nil {
				log.Println("service error")
				time.Sleep(time.Duration(retryInterval) * time.Second)
				continue
			}

			time.Sleep(time.Duration(sleepInterval) * time.Second)
		}
	}(videoService, youtubeService)

	server := NewServer(router, videoService)

	// start the server
	err = server.Run()

	if err != nil {
		return err
	}

	return nil
}
