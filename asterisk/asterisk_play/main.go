package main

import (
	"log"
	"time"

	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari/client/native" // Corrected import path
	"github.com/google/uuid"
)

const (
	ARIUsername  = "asterisk"                  // Define ARI username
	ARIPassword  = "asterisk"                  // Define ARI password
	AsteriskURL  = "http://localhost:8088/ari" // Default Asterisk ARI URL
	AppName      = "test"
	Endpoint     = "PJSIP/1000"          // Default endpoint
	PlaybackFile = "sound:demo-congrats" // Default playback file
)

func main() {

	// Create a new ARI client using native.Connect with options
	opts := &native.Options{
		Application: AppName,
		URL:         AsteriskURL,
		Username:    ARIUsername,
		Password:    ARIPassword,
	}

	ariClient, err := native.Connect(opts)
	if err != nil {
		log.Fatalf("Failed to connect to ARI: %v", err)
	}

	// Create a new channel
	channel, err := ariClient.Channel().Originate(nil, ari.OriginateRequest{
		Endpoint: Endpoint,
		App:      AppName,
	})

	if err != nil {
		log.Fatalf("Failed to originate channel: %v", err)
	}

	defer func() {
		err := channel.Hangup()
		if err != nil {
			log.Printf("Error hanging up channel: %v", err)
		}
	}()

	log.Printf("Originated channel: %s, key: %v\n", channel.ID(), channel.Key())

	// Give channel some time to establish
	time.Sleep(3 * time.Second)

	// Play a sound
	playbackID := uuid.New().String()
	playback, err := channel.Play(playbackID, PlaybackFile)
	if err != nil {
		log.Printf("Error playing sound: %v", err)
	}
	log.Printf("Started playback with ID: %s on channel: %s, playback: %v\n", playbackID, channel.ID(), playback)

	time.Sleep(1 * time.Second) // Wait for playback to start

	// Play a sound
	playbackID2 := uuid.New().String()
	playback2, err := channel.Play(playbackID2, PlaybackFile)
	if err != nil {
		log.Printf("Error playing sound: %v", err)
	}
	log.Printf("Started second playback with ID: %s on channel: %s, playback2: %v\n", playbackID2, channel.ID(), playback2)

	if playback != nil {
		defer func() {
			if err := playback.Stop(); err != nil {
				log.Printf("Error stopping playback: %v", err)
			}
		}()
	}

	// Keep channel alive for a few seconds
	time.Sleep(50 * time.Second)

	log.Println("Done.")
}
