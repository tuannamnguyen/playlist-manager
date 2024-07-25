package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dotenv-org/godotenvvault"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tuannamnguyen/playlist-manager/internal/repository"
	"github.com/tuannamnguyen/playlist-manager/internal/rest"
	"github.com/tuannamnguyen/playlist-manager/internal/service"
)

func main() {
	// setup .env
	err := godotenvvault.Load()
	if err != nil {
		log.Fatalf("error reading .env: %v", err)
	}

	// setup DB
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"playlist_manager",
	)
	db, err := sqlx.Connect("pgx", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// setup server
	e := echo.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()

	go startServer(e, db)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func startServer(e *echo.Echo, db *sqlx.DB) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s, World!", os.Getenv("HELLO")))
	})

	setupAPIRouter(e, db)

	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
}

func setupAPIRouter(e *echo.Echo, db *sqlx.DB) {
	apiRouter := e.Group("/api")

	apiRouter.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "You have been authenticated")
	})

	// setup playlist endpoint
	playlistRepository := repository.NewPlaylistRepository(db)
	songRepository := repository.NewSongRepository(db)
	playlistSongRepository := repository.NewPlaylistSongRepository(db)

	playlistService := service.NewPlaylist(playlistRepository, songRepository, playlistSongRepository)
	playlistHandler := rest.NewPlaylistHandler(playlistService)

	playlistRouter := apiRouter.Group("/playlists")

	playlistRouter.POST("", playlistHandler.Add)
	playlistRouter.GET("", playlistHandler.GetAll)
	playlistRouter.GET("/:id", playlistHandler.GetByID)
	playlistRouter.DELETE("/:id", playlistHandler.DeleteByID)

	// playlist-songs table endpoint
	playlistRouter.POST("/:playlist_id/songs", playlistHandler.AddSongsToPlaylist)
}
