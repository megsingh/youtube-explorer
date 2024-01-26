package main

import (
	"log"
	"os"
	"youtube_project/internal/app"

	"github.com/joho/godotenv"
)

func init() {
	loadEnvVariables()
}

func loadEnvVariables() {
	// Load environment variables from the .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

}

func main() {

	if err := app.Run(); err != nil {
		log.Fatalf("this is the startup error: %s\\n", err)
		os.Exit(1)
	}
}
