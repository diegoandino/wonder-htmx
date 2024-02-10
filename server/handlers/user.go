package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"text/template"

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
	LoginHandler  LoginHandler
}

func (h *UserHandler) SearchUsersHandler(c echo.Context) error {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err)
	}
	defer db.Close()

	query := c.QueryParam("query")
	fmt.Println("Query: ", query)

	if query == "" {
		return c.HTML(http.StatusOK, "")
	}

	stmt, err := db.Prepare(`select spotify_user_id, display_name, profile_picture from users 
				where display_name like ?
				and spotify_user_id != ?`)
	if err != nil {
		log.Fatal("Couldn't prepare db statement:", err)
	}
	defer stmt.Close()

	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(query, currentUser.ID)
	if err != nil {
		log.Fatal("Couldn't execute db query:", err)
	}
	defer rows.Close()

	var users []model.UserPayload
	for rows.Next() {
		var user model.UserPayload
		err := rows.Scan(&user.ID, &user.Username, &user.ProfilePicture)
		if err != nil {
			return err
		}

		fmt.Println("user: ", user)
		users = append(users, user)
	}

	// Prepare the partial HTML snippet for search results
	tmpl := template.New("searchResults")
	tmpl, err = tmpl.Parse(`
        <ul hx-swap-oob="true" id="search-results-dropdown">
            {{range .}}
                <li>
                    <h3>{{.Username}}</h3>
                    <img src="{{.ProfilePicture}}" alt="Profile Picture" style="width: 50px; height: 50px;">
					<button>Add Friend</button>
                </li>
            {{end}}
        </ul>
    `)

	if err != nil {
		return err
	}

	// Execute the template with the users slice to generate the HTML
	return tmpl.Execute(c.Response().Writer, users)
}

func (h UserHandler) UserShowHandler(c echo.Context) error {
	// Check if logged in
	sessionIsNew, err := h.LoginHandler.SessionExists(c)
	if err != nil {
		return err
	}

	// If it's a new session, redirect to login
	if sessionIsNew {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return nil
	}

	userPayload, err := h.getUserPayload(c)
	if err != nil {
		return err
	}

	dbErr := h.upsertUserDB(c, userPayload)
	if dbErr != nil {
		return dbErr
	}

	friends, friendErr := h.getFriends(c)
	if friendErr != nil {
		return err
	}

	return render(c, user.Show(userPayload, friends))
}

func (h UserHandler) upsertUserDB(c echo.Context, userPayload model.UserPayload) error {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err)
	}
	defer db.Close()

	// Prepare the UPSERT statement
	stmt, err := db.Prepare(`
        INSERT INTO Users(spotify_user_id, display_name, profile_picture, current_playing_song, current_album_art, current_album_name, current_artist_name, current_song_url) 
        VALUES(?,?,?,?,?,?,?,?)
        ON CONFLICT(spotify_user_id) 
        DO UPDATE SET 
            display_name=excluded.display_name, 
            profile_picture=excluded.profile_picture, 
            current_playing_song=excluded.current_playing_song, 
            current_album_art=excluded.current_album_art, 
            current_album_name=excluded.current_album_name, 
            current_artist_name=excluded.current_artist_name,
			current_song_url=excluded.current_song_url
    `)
	if err != nil {
		log.Fatal("Couldn't prepare db statement: upsertUser. Err: ", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(userPayload.ID, userPayload.Username, userPayload.ProfilePicture, userPayload.CurrentSongName, userPayload.CurrentAlbumArt, userPayload.CurrentAlbumName, userPayload.CurrentArtistName, userPayload.CurrentSongUrl)

	return err
}

func (h UserHandler) getFriends(c echo.Context) ([]model.UserPayload, error) {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err)
	}
	defer db.Close()

	var friends []model.UserPayload
	user, err := h.getCurrentUser(c)
	if err != nil {
		return nil, err
	}

	// Prepare the SELECT statement
	stmt, err := db.Prepare(`
        SELECT user_id_2 FROM FRIENDS 
		WHERE user_id_1=?
    `)
	if err != nil {
		log.Fatal("Couldn't prepare db statement:", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(user.ID)
	if err != nil {
		log.Fatal("Couldn't execute db query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var friendID string
		if err := rows.Scan(&friendID); err != nil {
			log.Fatal("Couldn't scan row:", err)
		}

		friendPayload, err := h.getFriendPayload(friendID)
		if err != nil {
			log.Printf("Couldn't get friend with ID %s: %v", friendID, err)
			return nil, err
		}

		friends = append(friends, friendPayload)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating rows:", err)
		return nil, err
	}

	return friends, nil
}

func (h UserHandler) getFriendPayload(friendID string) (model.UserPayload, error) {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err) // Consider changing log.Fatal to a more graceful error handling
	}
	defer db.Close()

	// Prepare the SELECT statement
	stmt, err := db.Prepare(`
        SELECT spotify_user_id, display_name, profile_picture, current_playing_song, current_album_art, current_album_name, current_artist_name, current_song_url 
        FROM USERS WHERE spotify_user_id=?
    `)
	if err != nil {
		log.Fatal("Couldn't prepare db statement:", err) // Consider changing log.Fatal to a more graceful error handling
	}
	defer stmt.Close()

	var friendPayload model.UserPayload
	err = stmt.QueryRow(friendID).Scan(&friendPayload.ID, &friendPayload.Username, &friendPayload.ProfilePicture, &friendPayload.CurrentSongName, &friendPayload.CurrentAlbumArt, &friendPayload.CurrentAlbumName, &friendPayload.CurrentArtistName, &friendPayload.CurrentSongUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.UserPayload{}, errors.New("No friend found with the given ID")
		}
		log.Fatal("Couldn't execute db query:", err) // Consider changing log.Fatal to a more graceful error handling
	}

	return friendPayload, nil
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
				ID:                user.ID,
				Username:          user.DisplayName,
				ProfilePicture:    user.Images[1].URL,
				CurrentAlbumArt:   playing.Item.Album.Images[1].URL,
				CurrentSongName:   playing.Item.Name,
				CurrentAlbumName:  playing.Item.Album.Name,
				CurrentArtistName: playing.Item.Artists[0].Name,
				CurrentSongUrl:    "https://open.spotify.com/track/" + string(playing.Item.URI)[14:],
			}
			c.Request().Header.Set("Content-Type", "application/json")
		} else {
			return model.UserPayload{}, c.String(200, "No song is currently playing")
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
