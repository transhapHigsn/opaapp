package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	pid := os.Getpid()
	filename := os.Getenv("ENV_FILE")
	log.Printf("pid=%d level=info msg=ENV_FILE: %s", pid, filename)
	if filename != "" {
		err := godotenv.Load(filename)
		if err != nil {
			log.Fatalf("pid=%d level=danger msg=Error loading %s file.", pid, filename)
		}
	}

	_, ok := os.LookupEnv("OPAAPP_ENV")
	if !ok {
		log.Fatalf("pid=%d level=danger msg=Error loading OPAAPP_ENV value.", pid)
	}
	fiberApp()
}
