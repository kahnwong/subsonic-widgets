package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// 3 requests per 10 seconds max
	app.Use(limiter.New(limiter.Config{
		Expiration: 10 * time.Second,
		Max:        3,
	}))

	// route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// entrypoint
	log.Fatal(app.Listen(":3000"))
}
