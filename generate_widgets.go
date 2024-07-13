package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
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

// utils
func sanitizeString(s string) string {
	return strings.Replace(s, "&", "&amp;", -1)
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

		// init template
		tmpl, err := template.ParseFiles("templates/now-playing.svg")
		if err != nil {
			log.Println("Template file doesn't exist")
		}

		// render template
		var tpl bytes.Buffer
		err = tmpl.Execute(&tpl, nowPlayingInfo)
		if err != nil {
			log.Println("Error rendering now playing")
		}
		fmt.Println(tpl.String())

		return base64.StdEncoding.EncodeToString(tpl.Bytes())
	}

	return ""
}
