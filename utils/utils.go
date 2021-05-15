package utils

import (
	"log"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("info=debug func=%s took %s", name, elapsed)
}
