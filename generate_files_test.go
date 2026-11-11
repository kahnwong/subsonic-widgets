package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCommandGeneratesWidgets(t *testing.T) {
	for _, playing := range []bool{false, true} {
		t.Run(fmt.Sprintf("playing=%t", playing), func(t *testing.T) {
			albumRequests := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/rest/getAlbumList":
					albumRequests++
					fmt.Fprint(w, `{"subsonic-response":{"albumList":{"album":[{"album":"Test album","artist":"Test artist","coverArt":"cover"}]}}}`)
				case "/rest/getNowPlaying":
					entry := "[]"
					if playing {
						entry = `[{"title":"Test track","artist":"Test artist","coverArt":"cover"}]`
					}
					fmt.Fprintf(w, `{"subsonic-response":{"nowPlaying":{"entry":%s}}}`, entry)
				case "/rest/getCoverArt":
					fmt.Fprint(w, "cover bytes")
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()
			oldEndpoint, oldMode, oldAuth := subsonicApiEndpoint, mode, authParams
			t.Cleanup(func() { subsonicApiEndpoint, mode, authParams = oldEndpoint, oldMode, oldAuth })
			subsonicApiEndpoint, mode = server.URL, ""
			t.Chdir(t.TempDir())
			logger := slog.New(slog.DiscardHandler)
			for range 2 {
				if err := run(nil, logger); err != nil {
					t.Fatal(err)
				}
			}
			if albumRequests != 10 {
				t.Fatalf("album requests = %d, want 10", albumRequests)
			}
			entries, err := os.ReadDir("data")
			if err != nil || len(entries) != 6 {
				t.Fatalf("output entries = %v, error = %v", entries, err)
			}
			names := []string{"now-playing.svg"}
			for i := range 5 {
				names = append(names, fmt.Sprintf("random-album-%d.svg", i+1))
			}
			for _, name := range names {
				data, err := os.ReadFile(filepath.Join("data", name))
				if err != nil || !strings.Contains(string(data), "<svg") {
					t.Fatalf("%s is not SVG: %v", name, err)
				}
			}
		})
	}
}

func TestCommandErrors(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	for _, args := range [][]string{{"unknown"}, {"server", "extra"}} {
		if err := run(args, logger); err == nil {
			t.Fatalf("run(%v) succeeded", args)
		}
	}
	oldMode, oldAuth := mode, authParams
	t.Cleanup(func() { mode, authParams = oldMode, oldAuth })
	mode = ""
	if err := run([]string{"server"}, logger); err == nil || !strings.Contains(err.Error(), "MODE") {
		t.Fatalf("server error = %v, want MODE error", err)
	}
}

func TestGenerationErrors(t *testing.T) {
	if _, err := generateRandomAlbumWidgetBase64(RandomAlbum{}); err == nil {
		t.Fatal("empty album response succeeded")
	}
	if err := writeSVG(filepath.Join(t.TempDir(), "bad.svg"), "!"); err == nil {
		t.Fatal("invalid base64 succeeded")
	}
	path := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := generateWidgets(path); err == nil {
		t.Fatal("invalid output directory succeeded")
	}
}
