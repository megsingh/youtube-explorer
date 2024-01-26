package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"youtube_project/internal/app"

	"github.com/joho/godotenv"
)

func init() {
	loadEnvVariables()
}

func loadEnvVariables() {

	// Get the absolute path to the directory containing main.go
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Println("Error getting source file path.")
		return
	}

	// Get the directory of the source file
	sourceDir := filepath.Dir(filename)

	// Construct the path to the .env file
	envFilePath := filepath.Join(sourceDir, "../../internal/config", ".env")

	// Load environment variables from the .env file
	err := godotenv.Load(envFilePath)
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
