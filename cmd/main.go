package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/PedroRibeiro95/syla/internal/handler"
	"github.com/PedroRibeiro95/syla/pkg/provider/spotify"
)

func main() {

	logInit("debug", "json")

	log.Info("Logging has been initialized")
	spotifyProvider := spotify.New("xxx", "xxx", "http://localhost:8080/auth")
	log.Debug("Instantiated Spotify Provider")

	spotifyHandler := handler.New(spotifyProvider)
	spotifyAuthHandler := handler.SpotifyAuthHandler{}
	log.Debug("Instantiated Spotify Handlers")

	// Registers Spotify Authenticator handler
	http.Handle("/auth", &spotifyAuthHandler)
	log.Debug("Registered Spotify authentication handler")

	// Register API handlers
	http.HandleFunc("/api/spotify/falbums", spotifyHandler.GetFavoriteAlbumsAPI())
	http.HandleFunc("/api/spotify/fartists", spotifyHandler.GetFavoriteArtistsAPI())
	log.Debug("Registered API handlers")

	fmt.Printf("\n\nPlease click the following link to allow Syla to access your Spotify information:\n   %s\n\n", spotifyProvider.URL)

	log.Debug("Waiting for Spotify authentication callback")
	go func() {
		// Waits for the authentication callback...
		for spotifyAuthHandler.Request == nil {
		}

		log.Debug("Callback received")
		spotifyProvider.InstantiateClient(spotifyAuthHandler.Request)
		log.Debug("Spotify client instantiated!")
	}()

	// Listens indefinetly
	log.Debug("Listening indefinetly on 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func logInit(level, formatter string) error {
	switch strings.ToUpper(level) {
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	default:
		log.Warn("The specified log level is invalid")
		return errors.New("The specified log level is invalid")
	}

	switch strings.ToUpper(formatter) {
	case "JSON":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Warn("Using default formatter")
		log.SetFormatter(&log.TextFormatter{})
	}

	log.SetOutput(os.Stdout)

	return nil
}
