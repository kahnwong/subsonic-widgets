package main

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// 20 requests per 1 minute max
	app.Use(limiter.New(limiter.Config{
		Expiration: 1 * time.Minute,
		Max:        20,
	}))

	// route
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
	log.Fatal(app.Listen(":3000"))
}
