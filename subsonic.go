package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/google/go-querystring/query"
	"github.com/rs/zerolog/log"
)

// auth
type SubsonicAuth struct {
	Username string `url:"u"`
	Token    string `url:"t"`
	Salt     string `url:"s"`
	Version  string `url:"v"`
	Client   string `url:"c"`
	Format   string `url:"f"`
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

// request params
type RandomAlbumRequest struct {
	Type string `url:"type"`
	Size int64  `url:"size"`
}

type CoverRequest struct {
	ID   string `url:"id"`
	Size int64  `url:"size"`
}

// fetchers
func getNowPlaying() NowPlaying {
	var response NowPlaying
	err := requests.
		URL(subsonicApiEndpoint).
		Method(http.MethodGet).
		Path("rest/getNowPlaying").
		Params(authParams).
		ToJSON(&response).
		Fetch(context.Background())

	if err != nil {
		log.Error().Msg("Failed to get NowPlaying")
	}

	return response
}

func getRandomAlbum() RandomAlbum {
	randomAlbumEnv := RandomAlbumRequest{
		Type: "random",
		Size: 1,
	}
	randomAlbumParams, _ := query.Values(randomAlbumEnv) // fetch response

	var response RandomAlbum
	err := requests.
		URL(subsonicApiEndpoint).
		Method(http.MethodGet).
		Path("rest/getAlbumList").
		Params(authParams).
		Params(randomAlbumParams).
		ToJSON(&response).
		Fetch(context.Background())

	if err != nil {
		log.Error().Msg("Failed to get RandomAlbum")
	}

	return response
}

func getCoverBase64(coverID string) string {
	coverEnv := CoverRequest{
		ID:   coverID,
		Size: 48,
	}
	coverParams, _ := query.Values(coverEnv)

	var buffer bytes.Buffer
	err := requests.
		URL(subsonicApiEndpoint).
		Method(http.MethodGet).
		Path("rest/getCoverArt").
		Params(authParams).
		Params(coverParams).
		ToBytesBuffer(&buffer).
		Fetch(context.Background())

	if err != nil {
		log.Error().Msg("Failed to get Cover")
	}

	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}
