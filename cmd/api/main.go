package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dotenv-org/godotenvvault"
	"github.com/gorilla/sessions"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/spotify"
	"github.com/tuannamnguyen/playlist-manager/internal/repository"
	"github.com/tuannamnguyen/playlist-manager/internal/rest"
	"github.com/tuannamnguyen/playlist-manager/internal/service"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"gopkg.in/boj/redistore.v1"
)

func main() {
	// setup .env
	err := godotenvvault.Load()
	if err != nil {
		log.Fatalf("error reading .env: %v", err)
	}

	// setup HTTP client
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   time.Minute,
		Transport: transport,
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

	// setup session and OAuth2
	redisInfo := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	key := os.Getenv("SESSION_SECRET")
	store, err := redistore.NewRediStore(10, "tcp", redisInfo, os.Getenv("REDIS_PASSWORD"), []byte(key))
	if err != nil {
		log.Fatalf("Unable to connect to Redis for session store: %v", err)
	}
	defer store.Close()
	store.SetMaxAge(3600)

	gob.Register(goth.User{})

	gothic.Store = store
	goth.UseProviders(
		spotify.New(
			os.Getenv("SPOTIFY_ID"),
			os.Getenv("SPOTIFY_SECRET"),
			os.Getenv("SPOTIFY_REDIRECT_URL"),
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistReadPrivate,
		),
	)

	// setup server
	e := echo.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()

	go startServer(e, db, httpClient, store)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func startServer(e *echo.Echo, db *sqlx.DB, httpClient *http.Client, store sessions.Store) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s, World!", os.Getenv("HELLO")))
	})

	setupAPIRouter(e, db, httpClient, store)

	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
}

func setupAPIRouter(e *echo.Echo, db *sqlx.DB, httpClient *http.Client, store sessions.Store) {
	apiRouter := e.Group("/api")

	apiRouter.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "You have been authenticated")
	})
	playlistRouter := apiRouter.Group("/playlists")
	searchRouter := apiRouter.Group("/search")
	oauthRouter := apiRouter.Group("/oauth")

	setupPlaylistRoutes(playlistRouter, db, store)
	setupSearchRoutes(searchRouter, httpClient)
	setupOAuthRoutes(oauthRouter, store)
}

func setupPlaylistRoutes(router *echo.Group, db *sqlx.DB, store sessions.Store) {
	// setup playlist endpoint
	playlistRepository := repository.NewPlaylistRepository(db)
	songRepository := repository.NewSongRepository(db)
	playlistSongRepository := repository.NewPlaylistSongRepository(db)
	albumRepository := repository.NewAlbumRepository(db)
	artistRepository := repository.NewArtistRepository(db)
	artistSongRepository := repository.NewArtistSongRepository(db)
	artistAlbumRepository := repository.NewArtistAlbumRepository(db)

	playlistService := service.NewPlaylist(
		playlistRepository,
		songRepository,
		playlistSongRepository,
		albumRepository,
		artistRepository,
		artistSongRepository,
		artistAlbumRepository,
	)
	playlistHandler := rest.NewPlaylistHandler(playlistService, store)

	// playlist CRUD
	router.POST("", playlistHandler.Add)
	router.GET("", playlistHandler.GetAll)
	router.GET("/:id", playlistHandler.GetByID)
	router.DELETE("/:id", playlistHandler.DeleteByID)

	// playlist-songs table endpoints
	playlistSongsEndpoint := "/:playlist_id/songs"
	router.POST(playlistSongsEndpoint, playlistHandler.AddSongsToPlaylist)
	router.GET(playlistSongsEndpoint, playlistHandler.GetAllSongsFromPlaylist)
	router.DELETE(playlistSongsEndpoint, playlistHandler.DeleteSongsFromPlaylist)

	// conversion endpoints
	router.POST("/:playlist_id/convert/:provider", playlistHandler.ConvertHandler)
}

func setupSearchRoutes(router *echo.Group, httpClient *http.Client) {
	searchRepository := repository.NewSearchRepository(httpClient)

	searchService := service.NewSearch(searchRepository)
	searchHandler := rest.NewSearchHandler(searchService)

	router.POST("", searchHandler.SearchMusicData)
}

func setupOAuthRoutes(router *echo.Group, store sessions.Store) {
	oauthHandler := rest.NewOAuthHandler(store)

	router.GET("/:provider", oauthHandler.LoginHandler)
	router.GET("/callback/:provider", oauthHandler.CallbackHandler)
	router.GET("/logout/:provider", oauthHandler.LogoutHandler)
}
