package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	pasteDir = "pastes"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/paste", func(c echo.Context) error {
		text, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		id := uuid.NewString()[:5] // Generate a 5-character ID
		err = savePaste(id, text)
		if err != nil {
			return err
		}
		return c.JSON(200, map[string]string{"id": id})
	})

	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.File(filepath.Join(pasteDir, id))
	})

	e.Static("/css", "web/css")
	e.Static("/", "web")

	log.Fatal(e.Start(":1323"))
}

func savePaste(id string, text []byte) error {
	if err := os.MkdirAll(pasteDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(pasteDir, id), text, 0644)
}
