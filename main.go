package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/google/go-querystring/query"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	slogfiber "github.com/samber/slog-fiber"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

var (
	mode = os.Getenv("MODE")

	subsonicApiEndpoint = os.Getenv("SUBSONIC_API_ENDPOINT")
	authValues          = SubsonicAuth{
		Username: os.Getenv("USERNAME"),
		Token:    os.Getenv("TOKEN"),
		Salt:     os.Getenv("SALT"),
		Version:  "1.16.1",
		Client:   "subsonic-widgets",
		Format:   "json",
	}
	authParams url.Values
)

func returnSVGResponse(c fiber.Ctx, svg string) error {
	c.Set(fiber.HeaderCacheControl, "no-cache")

	data, err := base64.StdEncoding.DecodeString(svg)
	if err != nil {
		slog.Error("Error decoding SVG", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error decoding SVG")
	}

	c.Set(fiber.HeaderContentType, "image/svg+xml")
	return c.Status(fiber.StatusOK).Send(data)
}

func newLogger(logLevel string, pretty bool) (*slog.Logger, error) {
	level := zerolog.DebugLevel
	var levelErr error
	logLevel = strings.TrimSpace(logLevel)
	if logLevel != "" {
		parsedLevel, err := zerolog.ParseLevel(strings.ToLower(logLevel))
		if err != nil {
			levelErr = err
		} else {
			level = parsedLevel
		}
	}

	logger := zerolog.New(os.Stderr).Level(level).With().Timestamp().Logger()
	if pretty {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return slog.New(slogzerolog.Option{
		Level:  slogzerolog.ZeroLogLeveler{Logger: &logger},
		Logger: &logger,
	}.NewZerologHandler()), levelErr
}

func newApp(logger *slog.Logger) *fiber.App {
	app := fiber.New()
	app.Use(slogfiber.New(logger))
	app.Use(recover.New())
	app.Use(limiter.New(limiter.Config{
		Max:        60,
		Expiration: time.Minute,
	}))

	// routes
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Welcome to subsonic-widgets api")
	})

	// --- now playing --- //
	app.Get("/now-playing.svg", func(c fiber.Ctx) error {
		nowPlaying, err := getNowPlaying()
		if err != nil {
			slog.Error("Failed to get now playing", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching now playing data")
		}

		svg, err := generateNowPlayingWidgetBase64(nowPlaying)
		if err != nil {
			slog.Error("Failed to generate now playing widget", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).SendString("Error generating widget")
		}

		return returnSVGResponse(c, svg)
	})

	// --- random album --- //
	for i := range 5 {
		app.Get(fmt.Sprintf("/random-album-%v.svg", i+1), func(c fiber.Ctx) error {
			randomAlbum, err := getRandomAlbum()
			if err != nil {
				slog.Error("Failed to get random album", slog.Any("error", err))
				return c.Status(fiber.StatusInternalServerError).SendString("Error fetching random album data")
			}

			svg, err := generateRandomAlbumWidgetBase64(randomAlbum)
			if err != nil {
				slog.Error("Failed to generate random album widget", slog.Any("error", err))
				return c.Status(fiber.StatusInternalServerError).SendString("Error generating widget")
			}

			return returnSVGResponse(c, svg)
		})
	}

	return app
}

func main() {
	logger, err := newLogger(os.Getenv("LOG_LEVEL"), mode == "development")
	slog.SetDefault(logger)
	if err != nil {
		slog.Error("Invalid LOG_LEVEL", slog.String("level", os.Getenv("LOG_LEVEL")), slog.Any("error", err))
		os.Exit(1)
	}

	listenAddress := ""
	switch mode {
	case "production":
		listenAddress = ":3000"
	case "development":
		listenAddress = "localhost:3000"
	default:
		slog.Error("Listen address is not set")
		os.Exit(1)
	}

	authParams, err = query.Values(authValues)
	if err != nil {
		slog.Error("Failed to create auth parameters", slog.Any("error", err))
		os.Exit(1)
	}

	if err := newApp(logger).Listen(listenAddress, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
		slog.Error("Fiber app error", slog.Any("error", err))
		os.Exit(1)
	}
}
