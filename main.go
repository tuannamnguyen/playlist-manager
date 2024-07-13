package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dotenv-org/godotenvvault"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenvvault.Load()
	if err != nil {
		log.Fatal("error reading .env")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s, World!", os.Getenv("HELLO")))
	})
	e.Logger.Fatal(e.Start(":8080"))
}
