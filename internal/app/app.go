package app

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
	"youtube_project/internal/api"
	"youtube_project/internal/repository"
	app_errors "youtube_project/pkg/error_handler"
	api_errors "youtube_project/pkg/error_handler/api_errors"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Run initializes and starts the application
func Run() error {

	// setup database connection
	dbClient, dbCollection, err := repository.SetupDatabase()
	if err != nil {
		return app_errors.NewDatabaseConnError(err.Error())
	}

	// create Storage instance
	storage := repository.NewStorage(dbClient, dbCollection)

	// create router instance
	router := gin.Default()
	router.Use(cors.Default())

	// create Youtube Data Api service
	youtubeService := api.NewYouTubeAPIService()

	// create Video service instance
	videoService := api.NewVideoService(storage)

	// start a goroutine for the background job of fetching videos from youtube
	go func(videoService api.VideoService, youtubeService api.YouTubeAPIService) {
		for {
			responseVideos, err := videoService.FetchFromYoutube(youtubeService)
			sleepInterval, _ := strconv.Atoi(os.Getenv("API_SLEEP_INTERVAL"))

			if err != nil {
				apiErr := api_errors.NewYoutubeAPIError(err.Error())

				// handle quota exceeding error for Youtube API
				if reflect.TypeOf(err) == reflect.TypeOf(api_errors.QuotaExceedError{}) {
					log.Println(apiErr.Error())

					youtubeService.RenewServiceAPIKey()
					continue
				} else {
					log.Println(apiErr.Error())
					os.Exit(1)
				}
			}

			if err == nil {

				// insert the obtained videos from Youtube into the DB
				insertErr := videoService.InsertInDB(responseVideos)
				if insertErr != nil {
					databaseErr := app_errors.NewDatabaseInsertionError(insertErr.Error())
					log.Println(databaseErr.Error())
					os.Exit(1)
				}
			}

			time.Sleep(time.Duration(sleepInterval) * time.Second)
		}
	}(videoService, youtubeService)

	//create Server Instance
	server := NewServer(router, videoService)

	// start the server
	err = server.Run()
	if err != nil {
		return app_errors.NewServerError(err.Error())
	}

	return nil
}
