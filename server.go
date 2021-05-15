package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	filename := os.Getenv("ENV_FILE")
	if filename == "" {
		filename = ".env"
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading %s file", filename)
	}
	fiberApp()
}
