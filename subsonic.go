package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/google/go-querystring/query"
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

type RandomAlbum struct {
	SubsonicResponse struct {
		Status        string `json:"status"`
		Version       string `json:"version"`
		Type          string `json:"type"`
		ServerVersion string `json:"serverVersion"`
		OpenSubsonic  bool   `json:"openSubsonic"`
		AlbumList     struct {
			Album []struct {
				ID            string        `json:"id"`
				Parent        string        `json:"parent"`
				IsDir         bool          `json:"isDir"`
				Title         string        `json:"title"`
				Name          string        `json:"name"`
				Album         string        `json:"album"`
				Artist        string        `json:"artist"`
				Year          int           `json:"year"`
				Genre         string        `json:"genre"`
				CoverArt      string        `json:"coverArt"`
				Duration      int           `json:"duration"`
				Created       time.Time     `json:"created"`
				ArtistID      string        `json:"artistId"`
				SongCount     int           `json:"songCount"`
				IsVideo       bool          `json:"isVideo"`
				Bpm           int           `json:"bpm"`
				Comment       string        `json:"comment"`
				SortName      string        `json:"sortName"`
				MediaType     string        `json:"mediaType"`
				MusicBrainzID string        `json:"musicBrainzId"`
				Genres        []interface{} `json:"genres"`
				ReplayGain    struct {
				} `json:"replayGain"`
				ChannelCount int       `json:"channelCount"`
				SamplingRate int       `json:"samplingRate"`
				PlayCount    int       `json:"playCount,omitempty"`
				Played       time.Time `json:"played,omitempty"`
			} `json:"album"`
		} `json:"albumList"`
	} `json:"subsonic-response"`
}

type RandomAlbumEnv struct {
	Type string `url:"type"`
	Size int64  `url:"size"`
}

type CoverEnv struct {
	ID   string `url:"id"`
	Size int64  `url:"size"`
}

// fetchers
func getNowPlaying() NowPlaying {
	requestUrl := fmt.Sprintf("%s/rest/getNowPlaying?%s", subsonicApiEndpoint, authParams.Encode())

	var response NowPlaying
	err := requests.
		URL(requestUrl).
		ToJSON(&response).
		Fetch(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	return response
}

func getRandomAlbum() RandomAlbum {
	randomAlbumEnv := RandomAlbumEnv{
		Type: "random",
		Size: 1,
	}
	randomAlbumParams, _ := query.Values(randomAlbumEnv) // fetch response

	requestUrl := fmt.Sprintf("%s/rest/getAlbumList?%s&%s", subsonicApiEndpoint, authParams.Encode(), randomAlbumParams.Encode())

	var response RandomAlbum
	err := requests.
		URL(requestUrl).
		ToJSON(&response).
		Fetch(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	return response
}

func getCoverBase64(coverID string) string {
	coverEnv := CoverEnv{
		ID:   coverID,
		Size: 48,
	}
	coverParams, _ := query.Values(coverEnv)

	requestUrl := fmt.Sprintf("%s/rest/getCoverArt?%s&%s", subsonicApiEndpoint, authParams.Encode(), coverParams.Encode())

	var buffer bytes.Buffer
	err := requests.
		URL(requestUrl).
		ToBytesBuffer(&buffer).
		Fetch(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}
