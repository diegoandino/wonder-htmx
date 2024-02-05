package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type UserHandler struct {
	Auth          spotify.Authenticator
	Store         *sessions.CookieStore
	UserDataStore *sync.Map
}

type UserPayload struct {
	Username         string `json:"username"`
	ProfilePicture   string `json:"profile_picture"`
	CurrentAlbumArt  string `json:"current_album_art"`
	CurrentSongName  string `json:"current_song_name"`
	CurrentAlbumName string `json:"current_album_name"`
}

func (h UserHandler) CurrentSongHandler(c echo.Context) error {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err)
	}
	defer db.Close()

	// Prepare the UPSERT statement
	stmt, err := db.Prepare(`
        INSERT INTO Users(spotify_user_id, display_name, current_playing_song, current_album_art, current_album_name, current_artist_name) 
        VALUES(?,?,?,?,?,?)
        ON CONFLICT(spotify_user_id) 
        DO UPDATE SET 
            display_name=excluded.display_name, 
            current_playing_song=excluded.current_playing_song, 
            current_album_art=excluded.current_album_art, 
            current_album_name=excluded.current_album_name, 
            current_artist_name=excluded.current_artist_name
    `)
	if err != nil {
		log.Fatal("Couldn't prepare db statement:", err)
	}
	defer stmt.Close()

	var userPayload UserPayload
	user, err := h.getCurrentUser(c)
	if err != nil {
		return err
	}

	if storedClient, ok := h.UserDataStore.Load(user.ID); ok {
		client := storedClient.(*spotify.Client)
		playing, err := client.PlayerCurrentlyPlaying()
		if err != nil {
			return err
		}

		if playing.Item != nil {
			//fmt.Fprintf(w, "Currently playing: %s by %s", playing.Item.Name, playing.Item.Artists[0].Name)
			_, err := stmt.Exec(user.ID, user.DisplayName, playing.Item.Name, playing.Item.Album.Images[0].URL, playing.Item.Album.Name, playing.Item.Artists[0].Name)
			if err != nil {
				log.Fatal("Couldn't upsert into Users table:", err)
			}

			userPayload = UserPayload{
				Username:         user.DisplayName,
				ProfilePicture:   user.Images[0].URL,
				CurrentAlbumArt:  playing.Item.Album.Images[0].URL,
				CurrentSongName:  playing.Item.Name,
				CurrentAlbumName: playing.Item.Album.Name,
			}
			c.Request().Header.Set("Content-Type", "application/json")
		} else {
			c.String(200, "No song is currently playing")
		}
	} else {
		c.String(404, "User data not found")
	}

	return c.JSON(http.StatusOK, userPayload)
}

func (h UserHandler) getCurrentUser(c echo.Context) (*spotify.PrivateUser, error) {
	session, _ := h.Store.Get(c.Request(), "spotify-session")
	token, ok := session.Values["token"].(*oauth2.Token)
	if !ok {
		return &spotify.PrivateUser{}, errors.New("Not Logged In")
	}

	client := h.Auth.NewClient(token)
	user, err := client.CurrentUser()
	if err != nil {
		return &spotify.PrivateUser{}, err
	}

	return user, nil
}
