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
	auth = spotify.NewAuthenticator(
		"https://wonder.andino.io/callback",
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopeUserReadPlaybackState,
	)
	store = sessions.NewCookieStore([]byte(loadEnv("STORE_KEY")))
	userDataStore sync.Map
)

func init() {
	gob.Register(&oauth2.Token{})
}

func loadEnv(key string) string {
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
	app.GET("/get-friends", userHandler.GetFriendsHandler)
	app.GET("/get-user-payload", userHandler.GetUserPayloadHandler)
	app.GET("/notifications", userHandler.LoadNotificationsHandler)
	app.GET("/check-notifications", userHandler.CheckNotificationsHandler)

	app.POST("/send-friend-request", userHandler.SendFriendRequestHandler)
	app.POST("/accept-friend-request", userHandler.AcceptFriendRequestHandler)
	app.POST("/decline-friend-request", userHandler.DeclineFriendRequestHandler)
	app.POST("/remove-friend", userHandler.RemoveFriendHandler)

	log.Println("Starting server on :8080")
	app.Logger.Fatal(app.Start(":8080"))
}
