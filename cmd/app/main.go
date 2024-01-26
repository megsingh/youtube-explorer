package main

import (
	"log"
	"os"
	"youtube_project/internal/app"
	app_errors "youtube_project/pkg/error_handler"

	"github.com/joho/godotenv"
)

func init() {
	loadEnvVariables()
}

func loadEnvVariables() {

	// Load environment variables from the .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		err = app_errors.NewEnvironmentVariableError("")
		log.Println(err.Error())
		os.Exit(1)
	}
}

func main() {

	if err := app.Run(); err != nil {
		err = app_errors.NewStartupError(err.Error())
		log.Println(err)
		os.Exit(1)
	}
}
