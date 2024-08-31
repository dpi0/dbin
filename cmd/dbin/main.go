package main

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	pasteDir = "pastes"
)

func zerologMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		stop := time.Now()

		log.Info().
			Int("status", c.Response().Status).
			Str("method", c.Request().Method).
			Str("uri", c.Request().URL.Path).
			Str("remote_ip", c.RealIP()).
			Str("user_agent", c.Request().UserAgent()).
			Dur("latency", stop.Sub(start)).
			Int64("bytes_in", c.Request().ContentLength).
			Int64("bytes_out", c.Response().Size).
			Msg("Request completed")

		return err
	}
}

func main() {
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	e := echo.New()
	e.Use(zerologMiddleware)
	e.Use(middleware.Recover())

	e.POST("/paste", func(c echo.Context) error {
		text, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Error().Err(err).Msg("Error reading request body")
			return err
		}
		id := uuid.NewString()[:6] // Generate a 5-character ID
		err = savePaste(id, text)
		if err != nil {
			log.Error().Err(err).Msg("Error saving paste")
			return err
		}
		log.Info().Str("id", id).Msg("Paste saved successfully")
		return c.JSON(200, map[string]string{"id": id})
	})

	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		pastePath := filepath.Join(pasteDir, id)
		if _, err := os.Stat(pastePath); os.IsNotExist(err) {
			log.Error().Str("id", id).Msg("Paste not found")
			return c.NoContent(404)
		}
		log.Info().Msgf("Serving paste /%s", id)
		return c.File(pastePath)
	})

	e.Static("/css", "web/css")
	e.Static("/", "web")

	log.Info().Msg("Server starting on port 1323")
	log.Fatal().Err(e.Start(":1323")).Msg("Error starting server")
}

func savePaste(id string, text []byte) error {
	if err := os.MkdirAll(pasteDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(pasteDir, id), text, 0644)
}
