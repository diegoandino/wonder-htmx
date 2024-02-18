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
	"github.com/diegoandino/wonder-go/views/notification"
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

func (h *UserHandler) LoadNotificationsHandler(c echo.Context) error {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err)
	}
	defer db.Close()

	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`select primary_id from friend_status where secondary_id = ? and status = 'pending'`)
	if err != nil {
		log.Fatal("Couldn't prepare db statement:", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(currentUser.ID)

	var ids []string
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return err
		}

		ids = append(ids, id)
	}

	var pendingRequestsSlice []model.UserPayload
	for _, id := range ids {
		var userPayload model.UserPayload
		stmt, err := db.Prepare(`select spotify_user_id, 
										display_name, 
										profile_picture, 
										current_playing_song, 
										current_album_art, 
										current_album_name, 
										current_artist_name, 
										current_song_url 
								from users where spotify_user_id = ?`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		err = stmt.QueryRow(id).Scan(&userPayload.ID, &userPayload.Username, &userPayload.ProfilePicture, &userPayload.CurrentSongName, &userPayload.CurrentAlbumArt, &userPayload.CurrentAlbumName, &userPayload.CurrentArtistName, &userPayload.CurrentSongUrl)
		if err != nil {
			return err
		}

		pendingRequestsSlice = append(pendingRequestsSlice, userPayload)
	}

	return render(c, notification.Show(pendingRequestsSlice))
}

func (h *UserHandler) RemoveFriendHandler(c echo.Context) error {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err)
	}
	defer db.Close()

	stmt, err := db.Prepare(`delete from friends where (user_id_1=? and user_id_2=?) or (user_id_1=? and user_id_2=?)`)
	if err != nil {
		log.Fatal("Couldn't prepare db statement:", err)
	}
	defer stmt.Close()

	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		return err
	}

	// Retrieve the secondary_user_id from the request
	secondaryUserID := c.FormValue("secondary_user_id")

	_, err = stmt.Exec(currentUser.ID, secondaryUserID, secondaryUserID, currentUser.ID)

	return err
}

func (h *UserHandler) SendFriendRequestHandler(c echo.Context) error {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err)
	}
	defer db.Close()

	stmt, err := db.Prepare(`insert into friend_status(primary_id, secondary_id, status) values(?,?,?)`)
	if err != nil {
		log.Fatal("Couldn't prepare db statement:", err)
	}
	defer stmt.Close()

	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		return err
	}

	// Retrieve the secondary_user_id from the request
	secondaryUserID := c.FormValue("secondary_user_id")

	_, err = stmt.Exec(currentUser.ID, secondaryUserID, "pending")

	sentFriendRequestAction := `
        <script>
            const addFriendBtn = document.getElementById('btn-add-friend');
            addFriendBtn.textContent = 'Friend request sent successfully.';
        </script>
    `
	return c.HTML(http.StatusOK, sentFriendRequestAction)
}

func (h *UserHandler) AcceptFriendRequestHandler(c echo.Context) error {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err)
	}
	defer db.Close()

	insertIntoFriendsStmt, err := db.Prepare(`insert into friends(user_id_1, user_id_2) values (?, ?)`)
	if err != nil {
		log.Fatal("Couldn't prepare insert into friends db statement:", err)
		return err
	}
	defer insertIntoFriendsStmt.Close()

	updateFriendStatusStmt, err := db.Prepare(`update friend_status set status='accepted' where 
											 (primary_id=? and secondary_id=?) or (primary_id=? and secondary_id=?)`)
	if err != nil {
		log.Fatal("Couldn't prepare update friend_status table db statement:", err)
		return err
	}
	defer updateFriendStatusStmt.Close()

	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		return err
	}

	secondaryUserID := c.FormValue("secondary_user_id")

	_, insertErr := insertIntoFriendsStmt.Exec(currentUser.ID, secondaryUserID)
	if insertErr != nil {
		return insertErr
	}

	_, updateErr := updateFriendStatusStmt.Exec(currentUser.ID, secondaryUserID, secondaryUserID, currentUser.ID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func (h *UserHandler) DeclineFriendRequestHandler(c echo.Context) error {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err)
	}
	defer db.Close()

	updateFriendStatusStmt, err := db.Prepare(`update friend_status set status='declined' where 
											 (primary_id=? and secondary_id=?) or (primary_id=? and secondary_id=?)`)
	if err != nil {
		log.Fatal("Couldn't prepare db statement:", err)
		return err
	}
	defer updateFriendStatusStmt.Close()

	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		return err
	}

	secondaryUserID := c.FormValue("secondary_user_id")
	_, err = updateFriendStatusStmt.Exec(currentUser.ID, secondaryUserID, secondaryUserID, currentUser.ID)

	return err
}

func (h *UserHandler) SearchUsersHandler(c echo.Context) error {
	db, err := sql.Open("sqlite3", "../db/wonder.db")
	if err != nil {
		log.Fatal("Couldn't open db:", err)
		return err
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
		return err
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

	alreadyFriendsStmt, err := db.Prepare(`
        SELECT user_id_1, user_id_2 FROM FRIENDS 
		WHERE user_id_1=? or user_id_2=?
    `)
	if err != nil {
		log.Fatal("Couldn't prepare db statement:", err)
		return err
	}
	defer alreadyFriendsStmt.Close()

	friendRows, err := alreadyFriendsStmt.Query(currentUser.ID, currentUser.ID)
	if err != nil {
		return err
	}
	defer friendRows.Close()

	for friendRows.Next() {
		var id1 string
		var id2 string
		err := friendRows.Scan(&id1, &id2)
		if err != nil {
			return err
		}

		for _, searchResultUserPayload := range users {
			if searchResultUserPayload.ID == id1 || searchResultUserPayload.ID == id2 {
				alreadyFriendsTmpl := template.New("searchResults")
				alreadyFriendsTmpl, err = alreadyFriendsTmpl.Parse(`
					<ul hx-swap-oob="true" id="search-results-dropdown" class="search-results flex flex-col p-4 md:p-0 mt-4 font-medium border border-gray-100 rounded-lg bg-gray-50 md:space-x-8 rtl:space-x-reverse md:flex-row md:mt-0 md:border-0 md:bg-white dark:bg-gray-800 md:dark:bg-gray-900 dark:border-gray-700">
						{{range .}}
							<li id="user-search-result" class="flex py-2 px-3 text-white bg-black rounded shadow-lg md:bg-transparent md:text-blue-700 md:p-0 md:dark:text-blue-500">
								<img src="{{.ProfilePicture}}" class="w-12 h-12 rounded-full mr-4" alt="Profile Picture" style="width: 50px; height: 50px;">
								<h3 class="nunito-bold mt-3 mr-3">{{.Username}}</h3>
								<button 
								id="btn-add-friend" 
								hx-post="/remove-friend" 
								hx-vals='{"secondary_user_id": "{{.ID}}"}' 
								hx-on:click="document.getElementById('user-search-result').style.display = 'none'; document.getElementById('search-results-dropdown').style.display = 'none';"
								class="bg-blue-500 nunito-bold text-sm hover:bg-blue-700 text-white py-2 px-4 rounded">
									Remove Friend
								</button>
							</li>
						{{end}}
					</ul>
				`)

				if err != nil {
					return err
				}

				return alreadyFriendsTmpl.Execute(c.Response().Writer, users)
			}
		}
	}

	addFriendTmpl := template.New("searchResults")
	addFriendTmpl, err = addFriendTmpl.Parse(`
        <ul hx-swap-oob="true" id="search-results-dropdown" class="flex flex-col p-4 md:p-0 mt-4 font-medium border border-gray-100 rounded-lg bg-gray-50 md:space-x-8 rtl:space-x-reverse md:flex-row md:mt-0 md:border-0 md:bg-white dark:bg-gray-800 md:dark:bg-gray-900 dark:border-gray-700">
            {{range .}}
                <li class="flex py-2 px-3 text-white bg-blue-700 rounded md:bg-transparent md:text-blue-700 md:p-0 md:dark:text-blue-500">
                    <img src="{{.ProfilePicture}}" class="w-12 h-12 rounded-full mr-4" alt="Profile Picture" style="width: 50px; height: 50px;">
                    <h3 class="nunito-medium mt-3 mr-3">{{.Username}}</h3>
					<button id="btn-add-friend" hx-post="/send-friend-request" hx-vals='{"secondary_user_id": "{{.ID}}"}' class="bg-blue-500 nunito-medium hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
					  Add Friend
					</button>
                </li>
            {{end}}
        </ul>
    `)

	if err != nil {
		return err
	}

	return addFriendTmpl.Execute(c.Response().Writer, users)
}

func (h UserHandler) UserShowHandler(c echo.Context) error {
	// Check if logged in
	token, valid, err := h.LoginHandler.ValidateSession(c)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// If it's a new session, redirect to login
	if !valid {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	// Use the valid token to create a new Spotify client
	client := h.LoginHandler.Auth.NewClient(token)

	// Store or update the client in UserDataStore
	currentUser, err := client.CurrentUser()
	if err != nil {
		log.Printf("Failed to get user info; Err: %v", err)
		return err
	}
	h.LoginHandler.UserDataStore.Store(currentUser.ID, &client)

	userPayload, err := h.getUserPayload(c, &client)
	if err != nil {
		fmt.Println(err)
		return err
	}

	dbErr := h.upsertUserDB(c, userPayload)
	if dbErr != nil {
		fmt.Println(err)
		return dbErr
	}

	friends, friendErr := h.getFriends(c)
	if friendErr != nil {
		fmt.Println(err)
		return err
	}

	return render(c, user.Show(userPayload, friends))
}

func (h UserHandler) GetFriendsHandler(c echo.Context) error {
	friends, err := h.getFriends(c)
	if err != nil {
		return err
	}

	return render(c, user.Friends(friends))
}

func (h UserHandler) GetUserPayloadHandler(c echo.Context) error {
	currentUser, err := h.getCurrentUser(c)
	if err != nil {
		return err
	}

	var client *spotify.Client
	if storedClient, ok := h.UserDataStore.Load(currentUser.ID); ok {
		client = storedClient.(*spotify.Client)
	} else {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	userPayload, err := h.getUserPayload(c, client)
	if err != nil {
		return err
	}

	dbErr := h.upsertUserDB(c, userPayload)
	if dbErr != nil {
		return dbErr
	}

	return render(c, user.CurrentUser(userPayload))
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

	stmt, err := db.Prepare(`
        SELECT user_id_1, user_id_2 FROM FRIENDS 
		WHERE user_id_1=? or user_id_2=?
    `)
	if err != nil {
		log.Fatal("Couldn't prepare db statement:", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(user.ID, user.ID)
	if err != nil {
		log.Fatal("Couldn't execute db query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id1 string
		var id2 string
		if err := rows.Scan(&id1, &id2); err != nil {
			log.Fatal("Couldn't scan row:", err)
		}

		// if id's aren't equal, it's the friend id
		// I know there's duped logic, will refactor later
		if id1 != user.ID {
			friendPayload, err := h.getFriendPayload(id1)
			if err != nil {
				log.Printf("Couldn't get friend with ID %s: %v", id1, err)
				return nil, err
			}

			friends = append(friends, friendPayload)
		} else if id2 != user.ID {
			friendPayload, err := h.getFriendPayload(id2)
			if err != nil {
				log.Printf("Couldn't get friend with ID %s: %v", id2, err)
				return nil, err
			}

			friends = append(friends, friendPayload)
		}
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

func (h UserHandler) getUserPayload(c echo.Context, client *spotify.Client) (model.UserPayload, error) {
	var userPayload model.UserPayload
	user, err := h.getCurrentUser(c)
	if err != nil {
		return model.UserPayload{}, err
	}

	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return model.UserPayload{}, err
	}

	if playing.Item != nil {
		userPayload = model.UserPayload{
			ID:                user.ID,
			Username:          user.DisplayName,
			ProfilePicture:    "https://ui-avatars.com/api/?background=fff&color=000&name=u",
			CurrentAlbumArt:   playing.Item.Album.Images[0].URL,
			CurrentSongName:   playing.Item.Name,
			CurrentAlbumName:  playing.Item.Album.Name,
			CurrentArtistName: playing.Item.Artists[0].Name,
			CurrentSongUrl:    "https://open.spotify.com/track/" + string(playing.Item.URI)[14:],
		}
		if len(user.Images) > 0 {
			userPayload.ProfilePicture = user.Images[0].URL
		}
		c.Request().Header.Set("Content-Type", "application/json")
	} else {
		return model.UserPayload{}, c.String(200, "No song is currently playing")
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
