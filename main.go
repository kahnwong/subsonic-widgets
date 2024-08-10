package main

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/google/go-querystring/query"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	logger zerolog.Logger
)

func returnSVGResponse(c *fiber.Ctx, svg string) error {
	c.Response().Header.Add("Cache-Control", "no-cache")

	data, err := base64.StdEncoding.DecodeString(svg)
	if err != nil {
		logger.Error().Err(err).Msgf("Error decoding SVG")
	}

	_, err = c.Write(data)

	return err
}

func init() {
	authParams, _ = query.Values(authValues)
}

func main() {
	// entrypoint
	listenAddress := ""
	isPrettyLog := false
	switch mode {
	case "production":
		listenAddress = ":3000"
	case "development":
		listenAddress = "localhost:3000"
		isPrettyLog = true
	default:
		log.Fatal().Msg("Listen address is not set")
	}

	// app
	app := fiber.New()
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	if isPrettyLog {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	// 60 requests per 1 minute max
	app.Use(limiter.New(limiter.Config{
		Expiration: 1 * time.Minute,
		Max:        60,
	}))

	// routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to subsonic-widgets api")
	})

	// --- now playing --- //
	app.Get("/now-playing.svg", func(c *fiber.Ctx) error {
		c.Type("svg")

		nowPlaying := getNowPlaying()
		svg := generateNowPlayingWidgetBase64(nowPlaying)

		return returnSVGResponse(c, svg)
	})

	// --- random album --- //
	for i := range 5 {
		app.Get(fmt.Sprintf("/random-album-%v.svg", i+1), func(c *fiber.Ctx) error {
			c.Type("svg")

			randomAlbum := getRandomAlbum()
			svg := generateRandomAlbumWidgetBase64(randomAlbum)

			return returnSVGResponse(c, svg)
		})
	}

	if err := app.Listen(listenAddress); err != nil {
		logger.Fatal().Err(err).Msg("Fiber app error")
	}
}
