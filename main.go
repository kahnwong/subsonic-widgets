package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-querystring/query"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
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
)

func returnSVGResponse(c *gin.Context, svg string) {
	c.Header("Cache-Control", "no-cache")

	data, err := base64.StdEncoding.DecodeString(svg)
	if err != nil {
		slog.Error("Error decoding SVG", slog.Any("error", err))
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

		slog.Info("request",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Duration("latency", latency),
			slog.String("ip", clientIP),
		)
	}
}

func main() {
	// entrypoint
	listenAddress := ""
	isPrettyLog := mode == "development"

	// logger setup
	level := zerolog.InfoLevel
	logLevel := strings.TrimSpace(os.Getenv("LOG_LEVEL"))
	var levelErr error
	if logLevel != "" {
		parsedLevel, err := zerolog.ParseLevel(strings.ToLower(logLevel))
		if err != nil {
			levelErr = err
		} else {
			level = parsedLevel
		}
	}

	logger := zerolog.New(os.Stderr).Level(level).With().Timestamp().Logger()
	if isPrettyLog {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	slog.SetDefault(slog.New(slogzerolog.Option{
		Level:  slogzerolog.ZeroLogLeveler{Logger: &logger},
		Logger: &logger,
	}.NewZerologHandler()))
	if levelErr != nil {
		slog.Error("Invalid LOG_LEVEL", slog.String("level", logLevel), slog.Any("error", levelErr))
		os.Exit(1)
	}

	switch mode {
	case "production":
		listenAddress = ":3000"
		gin.SetMode(gin.ReleaseMode)
	case "development":
		listenAddress = "localhost:3000"
		gin.SetMode(gin.DebugMode)
	default:
		slog.Error("Listen address is not set")
		os.Exit(1)
	}

	var err error
	authParams, err = query.Values(authValues)
	if err != nil {
		slog.Error("Failed to create auth parameters", slog.Any("error", err))
		os.Exit(1)
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
			slog.Error("Failed to get now playing", slog.Any("error", err))
			c.String(http.StatusInternalServerError, "Error fetching now playing data")
			return
		}

		svg, err := generateNowPlayingWidgetBase64(nowPlaying)
		if err != nil {
			slog.Error("Failed to generate now playing widget", slog.Any("error", err))
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
				slog.Error("Failed to get random album", slog.Any("error", err))
				c.String(http.StatusInternalServerError, "Error fetching random album data")
				return
			}

			svg, err := generateRandomAlbumWidgetBase64(randomAlbum)
			if err != nil {
				slog.Error("Failed to generate random album widget", slog.Any("error", err))
				c.String(http.StatusInternalServerError, "Error generating widget")
				return
			}

			returnSVGResponse(c, svg)
		})
	}

	if err := router.Run(listenAddress); err != nil {
		slog.Error("Gin app error", slog.Any("error", err))
		os.Exit(1)
	}
}
