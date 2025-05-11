package spotify

import (
	"log"
	"time"
)

type SystemController interface {
	GetCurrentPlayingTrackTitle() (string, error)
	OpenURLInBrowser(url string) error
}

type genius interface {
	GetTrackPageURL(title string) (string, error)
}

type Spotify struct {
	systemController SystemController
	genius           genius

	prevTitle string
}

func NewSpotify(systemController SystemController, genius genius) *Spotify {
	return &Spotify{
		systemController: systemController,
		genius:           genius,
		prevTitle:        "",
	}
}

func (s *Spotify) Run(scanInterval time.Duration) {
	log.Println("Starting scanning Spotify...")

	for {
		title, err := s.systemController.GetCurrentPlayingTrackTitle()
		if err != nil {
			log.Println("Spotify is not running...Waiting 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		if title != s.prevTitle {
			log.Println("New track:", title)
			s.prevTitle = title

			url, err := s.genius.GetTrackPageURL(title)
			if err != nil {
				log.Printf("Failed to open Genius page: %v", err)
			}

			if err := s.systemController.OpenURLInBrowser(url); err != nil {
				log.Printf("Failed to open Genius page: %v", err)
			}
		}

		time.Sleep(scanInterval)
	}
}
