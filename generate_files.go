package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

func generateWidgets(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	for i := range 5 {
		name := fmt.Sprintf("random-album-%d.svg", i+1)
		album, err := getRandomAlbum()
		if err != nil {
			return fmt.Errorf("fetch %s: %w", name, err)
		}
		svg, err := generateRandomAlbumWidgetBase64(album)
		if err != nil {
			return fmt.Errorf("generate %s: %w", name, err)
		}
		if err := writeSVG(filepath.Join(dir, name), svg); err != nil {
			return err
		}
	}

	nowPlaying, err := getNowPlaying()
	if err != nil {
		return fmt.Errorf("fetch now playing: %w", err)
	}
	svg, err := generateNowPlayingWidgetBase64(nowPlaying)
	if err != nil {
		return fmt.Errorf("generate now playing: %w", err)
	}
	return writeSVG(filepath.Join(dir, "now-playing.svg"), svg)
}

func writeSVG(path, encoded string) error {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
