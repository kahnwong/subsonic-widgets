package main

import (
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
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// entrypoint
	log.Fatal(app.Listen(":3000"))
}
