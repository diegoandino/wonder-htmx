package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/diegoandino/wonder-go/model"
	"github.com/diegoandino/wonder-go/views/user"
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

func (h UserHandler) UserShowHandler(c echo.Context) error {
	userPayload, err := h.getUserPayload(c)
	if err != nil {
		return err
	}
	return render(c, user.Show(userPayload))
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
	userPayload, err := h.getUserPayload(c)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userPayload.ID, userPayload.Username, userPayload.CurrentSongName, userPayload.CurrentAlbumArt, userPayload.CurrentAlbumName, userPayload.CurrentArtistName)

	return c.JSON(http.StatusOK, userPayload)
}

func (h UserHandler) getUserPayload(c echo.Context) (model.UserPayload, error) {
	var userPayload model.UserPayload
	user, err := h.getCurrentUser(c)
	if err != nil {
		return model.UserPayload{}, err
	}

	if storedClient, ok := h.UserDataStore.Load(user.ID); ok {
		client := storedClient.(*spotify.Client)
		playing, err := client.PlayerCurrentlyPlaying()
		if err != nil {
			return model.UserPayload{}, err
		}

		if playing.Item != nil {
			if err != nil {
				log.Fatal("Couldn't upsert into Users table:", err)
			}

			userPayload = model.UserPayload{
				ID:               user.ID,
				Username:         user.DisplayName,
				ProfilePicture:   user.Images[0].URL,
				CurrentAlbumArt:  playing.Item.Album.Images[0].URL,
				CurrentSongName:  playing.Item.Name,
				CurrentAlbumName: playing.Item.Album.Name,
			}
			c.Request().Header.Set("Content-Type", "application/json")
		} else {
			return model.UserPayload{}, errors.New("No song is currently playing")
		}
	} else {
		return model.UserPayload{}, errors.New("User data not found")
	}

	return userPayload, nil
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
