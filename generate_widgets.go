package main

import (
	"bytes"
	"embed"
	"encoding/base64"
	"strings"
	"text/template"
)

//go:embed templates/*
var templatesFS embed.FS

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

func renderTemplateBase64(templatePath string, data any) (string, error) {
	// init template
	tmpl, err := template.ParseFS(templatesFS, templatePath)
	if err != nil {
		return "", err
	}

	// render template
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(tpl.Bytes()), nil
}

// generators
func generateNowPlayingWidgetBase64(nowPlaying NowPlaying) (string, error) {
	track := nowPlaying.SubsonicResponse.NowPlaying.Entry
	if len(track) == 1 {
		coverBase64, err := getCoverBase64(track[0].CoverArt)
		if err != nil {
			return "", err
		}

		nowPlayingInfo := NowPlayingInfo{
			Title:       sanitizeString(track[0].Title),
			Artist:      sanitizeString(track[0].Artist),
			CoverBase64: coverBase64,
		}

		return renderTemplateBase64("templates/now-playing.svg", nowPlayingInfo)
	}

	return renderTemplateBase64("templates/now-playing-null.svg", nil)
}

func generateRandomAlbumWidgetBase64(randomAlbum RandomAlbum) (string, error) {
	album := randomAlbum.SubsonicResponse.AlbumList.Album[0]

	coverBase64, err := getCoverBase64(album.CoverArt)
	if err != nil {
		return "", err
	}

	albumInfo := RandomAlbumInfo{
		Album:       sanitizeString(album.Album),
		Artist:      sanitizeString(album.Artist),
		CoverBase64: coverBase64,
	}

	return renderTemplateBase64("templates/random-album.svg", albumInfo)
}
