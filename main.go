package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
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

func returnSVGResponse(c *gin.Context, svg string) {
	c.Header("Cache-Control", "no-cache")

	data, err := base64.StdEncoding.DecodeString(svg)
	if err != nil {
		logger.Error().Err(err).Msg("Error decoding SVG")
		c.String(http.StatusInternalServerError, "Error decoding SVG")
		return
	}

	c.Data(http.StatusOK, "image/svg+xml", data)
}

func zerologMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info().
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("ip", clientIP).
			Msg("request")
	}
}

func init() {
	var err error
	authParams, err = query.Values(authValues)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create auth parameters")
	}
}

func main() {
	// entrypoint
	listenAddress := ""
	isPrettyLog := false
	switch mode {
	case "production":
		listenAddress = ":3000"
		gin.SetMode(gin.ReleaseMode)
	case "development":
		listenAddress = "localhost:3000"
		isPrettyLog = true
		gin.SetMode(gin.DebugMode)
	default:
		log.Fatal().Msg("Listen address is not set")
	}

	// logger setup
	logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	if isPrettyLog {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// rate limiter setup - 60 requests per 1 minute max
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  60,
	}
	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate)

	// app
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(zerologMiddleware())
	router.Use(mgin.NewMiddleware(rateLimiter))

	// routes
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to subsonic-widgets api")
	})

	// --- now playing --- //
	router.GET("/now-playing.svg", func(c *gin.Context) {
		nowPlaying, err := getNowPlaying()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get now playing")
			c.String(http.StatusInternalServerError, "Error fetching now playing data")
			return
		}

		svg, err := generateNowPlayingWidgetBase64(nowPlaying)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to generate now playing widget")
			c.String(http.StatusInternalServerError, "Error generating widget")
			return
		}

		returnSVGResponse(c, svg)
	})

	// --- random album --- //
	for i := range 5 {
		router.GET(fmt.Sprintf("/random-album-%v.svg", i+1), func(c *gin.Context) {
			randomAlbum, err := getRandomAlbum()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get random album")
				c.String(http.StatusInternalServerError, "Error fetching random album data")
				return
			}

			svg, err := generateRandomAlbumWidgetBase64(randomAlbum)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to generate random album widget")
				c.String(http.StatusInternalServerError, "Error generating widget")
				return
			}

			returnSVGResponse(c, svg)
		})
	}

	if err := router.Run(listenAddress); err != nil {
		logger.Fatal().Err(err).Msg("Gin app error")
	}
}
