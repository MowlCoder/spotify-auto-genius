package main

import (
	"log"
	"time"

	"github.com/MowlCoder/spotify-auto-genius/internal/genius"
	"github.com/MowlCoder/spotify-auto-genius/internal/spotify"
	"github.com/MowlCoder/spotify-auto-genius/internal/system"
)

func getSystemController() (spotify.SystemController, error) {
	return system.NewSystemController()
}

func main() {
	systemController, err := getSystemController()
	if err != nil {
		log.Fatalf("Failed to get system controller: %v", err)
	}

	genius := genius.NewGenius()

	spotifyWorker := spotify.NewSpotify(systemController, genius)
	spotifyWorker.Run(1 * time.Second)
}
