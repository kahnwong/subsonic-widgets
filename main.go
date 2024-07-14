package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/google/go-querystring/query"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	_ "github.com/joho/godotenv/autoload"
)

var (
	subsonicApiEndpoint = os.Getenv("SUBSONIC_API_ENDPOINT")
	authParams          url.Values
)

func returnSVGResponse(c *fiber.Ctx, svg string) error {
	c.Response().Header.Add("Cache-Control", "no-cache")

	data, err := base64.StdEncoding.DecodeString(svg)
	if err != nil {
		log.Fatal("error:", err)
	}

	_, err = c.Write(data)

	return err
}

func main() {
	app := fiber.New()

	// 20 requests per 1 minute max
	app.Use(limiter.New(limiter.Config{
		Expiration: 1 * time.Minute,
		Max:        20,
	}))

	// init
	authParams, _ = query.Values(authEnv)

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

	// entrypoint
	mode := os.Getenv("MODE")
	listenAddress := ""
	if mode == "production" {
		listenAddress = ":3000"
	} else if mode == "development" {
		listenAddress = "localhost:3000"
	} else {
		fmt.Println("Listen address is not set")
	}
	log.Fatal(app.Listen(listenAddress))
}
