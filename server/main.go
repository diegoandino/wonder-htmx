package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
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
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/current-song", currentSongHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	state := "example-state" // Replace with a secure random state
	url := auth.AuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	_, code := r.FormValue("state"), r.FormValue("code")
	token, err := auth.Exchange(code)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	session, _ := store.Get(r, "spotify-session")
	session.Values["token"] = token
	err = session.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Couldn't save session", http.StatusInternalServerError)
		return
	}

	client := auth.NewClient(token)
	user, err := client.CurrentUser()
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	userDataStore.Store(user.ID, &client)
	fmt.Fprintf(w, "Login successful! User ID: %s", user.ID)
}

func currentSongHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(w, r)
	if err != nil {
		http.Error(w, "Couldn't get user id", http.StatusInternalServerError)
	}

	if storedClient, ok := userDataStore.Load(userId); ok {
		client := storedClient.(*spotify.Client)
		playing, err := client.PlayerCurrentlyPlaying()
		if err != nil {
			http.Error(w, "Failed to get current playing song", http.StatusInternalServerError)
			return
		}

		if playing.Item != nil {
			fmt.Fprintf(w, "Currently playing: %s by %s", playing.Item.Name, playing.Item.Artists[0].Name)
		} else {
			fmt.Fprintf(w, "No song is currently playing.")
		}
	} else {
		http.Error(w, "User data not found", http.StatusNotFound)
	}
}

func getUserId(w http.ResponseWriter, r *http.Request) (string, error) {
	session, _ := store.Get(r, "spotify-session")
	token, ok := session.Values["token"].(*oauth2.Token)
	if !ok {
		http.Error(w, "You're not logged in", http.StatusUnauthorized)
		return "", errors.New("Not Logged In")
	}

	client := auth.NewClient(token)
	user, err := client.CurrentUser()
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return "", err
	}

	return user.ID, nil
}
