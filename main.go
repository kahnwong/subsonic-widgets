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

func main() {
	app := fiber.New()

	// 20 requests per 1 minute max
	app.Use(limiter.New(limiter.Config{
		Expiration: 1 * time.Minute,
		Max:        20,
	}))

	// test
	authParams, _ = query.Values(authEnv)

	nowPlaying := getNowPlaying()
	generateNowPlayingWidget(nowPlaying)

	// routes
	app.Get("/image.svg", func(c *fiber.Ctx) error {
		c.Type("svg")

		s := "base64ImagePlaceholder"
		data, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			log.Fatal("error:", err)
		}

		_, err = c.Write(data)
		return err
	})

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
