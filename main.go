package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := fiber.New()

	// 20 requests per 1 minute max
	app.Use(limiter.New(limiter.Config{
		Expiration: 1 * time.Minute,
		Max:        20,
	}))

	// test
	p := getNowPlaying(authEnv)
	fmt.Println(p.SubsonicResponse.NowPlaying.Entry[0].Album)

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
