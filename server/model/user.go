package model

type UserPayload struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	ProfilePicture    string `json:"profile_picture"`
	CurrentAlbumArt   string `json:"current_album_art"`
	CurrentSongName   string `json:"current_song_name"`
	CurrentArtistName string `json:"current_artist_name"`
	CurrentAlbumName  string `json:"current_album_name"`
	CurrentSongUrl    string `json:"current_song_url"`
}

type UserPayloadTemplate struct {
	UserPayload
	AlreadyFriend bool
	Pending       bool
}
