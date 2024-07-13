package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"log"
)

// auth
type subsonicAuth struct {
	Username string `url:"u"`
	Token    string `url:"t"`
	Salt     string `url:"s"`
	Version  string `url:"v"`
	Client   string `url:"c"`
	Format   string `url:"f"`
}

var authEnv = subsonicAuth{
	Username: os.Getenv("USERNAME"),
	Token:    os.Getenv("TOKEN"),
	Salt:     os.Getenv("SALT"),
	Version:  "1.16.1",
	Client:   "subsonic-widgets",
	Format:   "json",
}

// structs
type NowPlaying struct {
	SubsonicResponse struct {
		Status        string `json:"status"`
		Version       string `json:"version"`
		Type          string `json:"type"`
		ServerVersion string `json:"serverVersion"`
		NowPlaying    struct {
			Entry []struct {
				ID          string    `json:"id"`
				Parent      string    `json:"parent"`
				IsDir       bool      `json:"isDir"`
				Title       string    `json:"title"`
				Album       string    `json:"album"`
				Artist      string    `json:"artist"`
				Track       int       `json:"track"`
				Year        int       `json:"year"`
				Genre       string    `json:"genre"`
				CoverArt    string    `json:"coverArt"`
				Size        int       `json:"size"`
				ContentType string    `json:"contentType"`
				Suffix      string    `json:"suffix"`
				Duration    int       `json:"duration"`
				BitRate     int       `json:"bitRate"`
				Path        string    `json:"path"`
				DiscNumber  int       `json:"discNumber"`
				Created     time.Time `json:"created"`
				AlbumID     string    `json:"albumId"`
				ArtistID    string    `json:"artistId"`
				Type        string    `json:"type"`
				IsVideo     bool      `json:"isVideo"`
				Username    string    `json:"username"`
				MinutesAgo  int       `json:"minutesAgo"`
				PlayerID    int       `json:"playerId"`
				PlayerName  string    `json:"playerName"`
			} `json:"entry"`
		} `json:"nowPlaying"`
	} `json:"subsonic-response"`
}

// fetchers
func getNowPlaying(authParams url.Values) NowPlaying {
	log.Println("Fetching now playing")

	// fetch response
	requestUrl := fmt.Sprintf("%s/rest/getNowPlaying?%s", os.Getenv("SUBSONIC_API_ENDPOINT"), authParams.Encode())
	resp, err := http.Get(requestUrl)
	if err != nil {
		log.Println("No response from request")
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body")
	}

	var response NowPlaying
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println("Can not unmarshal JSON")
	}

	return response
}
