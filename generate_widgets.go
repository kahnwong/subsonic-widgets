package main

import (
	"bytes"
	"encoding/base64"
	"log"
	"strings"
	"text/template"
)

// structs
type NowPlayingInfo struct {
	Title       string
	Artist      string
	CoverBase64 string
}

type RandomAlbumInfo struct {
	Album       string
	Artist      string
	CoverBase64 string
}

// utils
func sanitizeString(s string) string {
	return strings.Replace(s, "&", "&amp;", -1)
}

func renderTemplateBase64(templatePath string, data any) string {
	// init template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("Error parsing template: %s", templatePath)
	}

	// render template
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		log.Println("Error rendering now playing")
	}

	return base64.StdEncoding.EncodeToString(tpl.Bytes())
}

// generators
func generateNowPlayingWidgetBase64(nowPlaying NowPlaying) string {
	track := nowPlaying.SubsonicResponse.NowPlaying.Entry
	if len(track) == 1 {
		nowPlayingInfo := NowPlayingInfo{
			Title:       sanitizeString(track[0].Title),
			Artist:      sanitizeString(track[0].Artist),
			CoverBase64: getCoverBase64(track[0].CoverArt),
		}

		return renderTemplateBase64("templates/now-playing.svg", nowPlayingInfo)
	} else {
		return renderTemplateBase64("templates/now-playing-null.svg", nil)
	}
}

func generateRandomAlbumWidgetBase64(randomAlbum RandomAlbum) string {
	album := randomAlbum.SubsonicResponse.AlbumList.Album[0]

	albumInfo := RandomAlbumInfo{
		Album:       sanitizeString(album.Album),
		Artist:      sanitizeString(album.Artist),
		CoverBase64: getCoverBase64(album.CoverArt),
	}

	return renderTemplateBase64("templates/random-album.svg", albumInfo)
}
