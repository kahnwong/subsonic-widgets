package main

import (
	"fmt"
	"strings"
)

func generateNowPlayingWidget(nowPlaying NowPlaying) string {
	track := nowPlaying.SubsonicResponse.NowPlaying.Entry
	if len(track) == 1 {
		title := strings.Replace(track[0].Title, "&", "&amp;", -1)
		album := strings.Replace(track[0].Album, "&", "&amp;", -1)
		artist := strings.Replace(track[0].Artist, "&", "&amp;", -1)
		coverBase64 := getCoverBase64(track[0].CoverArt)

		fmt.Println(title, album, artist, coverBase64)
	}

	return ""
}
