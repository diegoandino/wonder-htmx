package main

import (
	"encoding/gob"
	"log"
	"os"
	"sync"

	"github.com/diegoandino/wonder-go/handlers"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	// Initialize the Spotify authenticator with your client ID and client secret
	auth = spotify.NewAuthenticator(
		"http://localhost:8080/callback",
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
	)

	// Create a sessions store (replace "your-secret-key" with a real secret key)
	store = sessions.NewCookieStore([]byte("your-secret-key"))

	// Use a concurrent map to store user-specific data
	userDataStore sync.Map
)

func init() {
	// Register the oauth2.Token type with the encoding/gob package
	gob.Register(&oauth2.Token{})
}

// use godot package to load/read the .env file and
// return the value of the key
func loadEnv(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	spotifyClientId := loadEnv("SPOTIFY_CLIENT_ID")
	spotifySecret := loadEnv("SPOTIFY_SECRET")
	auth.SetAuthInfo(spotifyClientId, spotifySecret)

	app := echo.New()

	app.Static("/static", "static")
	loginHandler := handlers.LoginHandler{Auth: auth, Store: store, UserDataStore: &userDataStore}
	userHandler := handlers.UserHandler{Auth: auth, Store: store, UserDataStore: &userDataStore, LoginHandler: loginHandler}
	app.GET("/", loginHandler.RedirectHandler)
	app.GET("/login", loginHandler.LoginHandler)
	app.GET("/callback", loginHandler.CallbackHandler)
	app.GET("/home", userHandler.UserShowHandler)
	app.GET("/search-friends", userHandler.SearchUsersHandler)

	log.Println("Starting server on :8080")
	app.Logger.Fatal(app.Start(":8080"))
}
