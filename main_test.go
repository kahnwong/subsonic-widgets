package main

import (
	"context"
	"encoding/base64"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestNewLoggerLevels(t *testing.T) {
	tests := []struct {
		name         string
		level        string
		debugEnabled bool
		infoEnabled  bool
		errorEnabled bool
	}{
		{name: "default", debugEnabled: true, infoEnabled: true, errorEnabled: true},
		{name: "info", level: " info ", infoEnabled: true, errorEnabled: true},
		{name: "error", level: "ERROR", errorEnabled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := newLogger(tt.level, false)
			if err != nil {
				t.Fatalf("newLogger() error = %v", err)
			}

			ctx := context.Background()
			if got := logger.Enabled(ctx, slog.LevelDebug); got != tt.debugEnabled {
				t.Errorf("debug enabled = %v, want %v", got, tt.debugEnabled)
			}
			if got := logger.Enabled(ctx, slog.LevelInfo); got != tt.infoEnabled {
				t.Errorf("info enabled = %v, want %v", got, tt.infoEnabled)
			}
			if got := logger.Enabled(ctx, slog.LevelError); got != tt.errorEnabled {
				t.Errorf("error enabled = %v, want %v", got, tt.errorEnabled)
			}
		})
	}
}

func TestNewLoggerRejectsInvalidLevel(t *testing.T) {
	if _, err := newLogger("verbose", false); err == nil {
		t.Fatal("newLogger() error = nil, want invalid level error")
	}
}

func TestRootRoute(t *testing.T) {
	app := newApp(slog.New(slog.DiscardHandler))
	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			t.Errorf("response.Body.Close() error = %v", err)
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if got, want := string(body), "Welcome to subsonic-widgets api"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestReturnSVGResponse(t *testing.T) {
	const svg = `<svg xmlns="http://www.w3.org/2000/svg"></svg>`
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		return returnSVGResponse(c, base64.StdEncoding.EncodeToString([]byte(svg)))
	})

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			t.Errorf("response.Body.Close() error = %v", err)
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}
	if got := response.Header.Get("Cache-Control"); got != "no-cache" {
		t.Errorf("Cache-Control = %q, want %q", got, "no-cache")
	}
	if got := response.Header.Get("Content-Type"); got != "image/svg+xml" {
		t.Errorf("Content-Type = %q, want %q", got, "image/svg+xml")
	}
	if got := string(body); got != svg {
		t.Errorf("body = %q, want %q", got, svg)
	}
}

func TestRateLimit(t *testing.T) {
	app := newApp(slog.New(slog.DiscardHandler))
	for requestNumber := 1; requestNumber <= 61; requestNumber++ {
		response, err := app.Test(httptest.NewRequest(http.MethodGet, "/", http.NoBody))
		if err != nil {
			t.Fatalf("request %d: app.Test() error = %v", requestNumber, err)
		}
		if err := response.Body.Close(); err != nil {
			t.Fatalf("request %d: response.Body.Close() error = %v", requestNumber, err)
		}

		want := http.StatusOK
		if requestNumber == 61 {
			want = http.StatusTooManyRequests
		}
		if response.StatusCode != want {
			t.Fatalf("request %d: status = %d, want %d", requestNumber, response.StatusCode, want)
		}
	}
}
