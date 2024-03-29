package handlers

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type LoginHandler struct {
	Auth          spotify.Authenticator
	Store         *sessions.CookieStore
	UserDataStore *sync.Map
}

func (h LoginHandler) RedirectHandler(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func (h LoginHandler) LoginHandler(c echo.Context) error {
	state := "example-state" // Replace with a secure random state
	url := h.Auth.AuthURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h LoginHandler) CallbackHandler(c echo.Context) error {
	_, code := c.QueryParam("state"), c.QueryParam("code")
	token, err := h.Auth.Exchange(code)
	if err != nil {
		return err
	}

	session, _ := h.Store.Get(c.Request(), "spotify-session")
	session.Values["token"] = token
	err = session.Save(c.Request(), c.Response().Writer)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		return err
	}

	client := h.Auth.NewClient(token)
	user, err := client.CurrentUser()
	if err != nil {
		log.Printf("Failed to get user info; Err: %v", err)
		return err
	}

	h.UserDataStore.Store(user.ID, &client)
	return c.Redirect(http.StatusTemporaryRedirect, "/home")
}

func (h LoginHandler) ValidateSession(c echo.Context) (*oauth2.Token, bool, error) {
	session, err := h.Store.Get(c.Request(), "spotify-session")
	if err != nil {
		return nil, false, err
	}

	token, ok := session.Values["token"].(*oauth2.Token)
	if !ok || token == nil {
		return nil, false, nil // No token found
	}

	// Check if the token is expired
	if token.Expiry.Before(time.Now()) {
		return token, false, nil // Token is expired
	}

	return token, true, nil
}
