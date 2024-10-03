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
	"github.com/go-playground/validator"
	"github.com/gorilla/sessions"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/spotify"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tuannamnguyen/playlist-manager/internal/repository"
	"github.com/tuannamnguyen/playlist-manager/internal/rest"
	"github.com/tuannamnguyen/playlist-manager/internal/rest/gothprovider/genius"
	"github.com/tuannamnguyen/playlist-manager/internal/service"
	"gopkg.in/boj/redistore.v1"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

var bucketName = "playlist-cover"

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

	// setup minio client
	endpoint := os.Getenv("OBJECT_STORAGE_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretKeyID := os.Getenv("MINIO_SECRET_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretKeyID, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("error setting up minio client: %s", err)
	}

	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalf("error checking bucket exists: %s", err)
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("error creating new bucket: %s", err)
		}
	}

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
			spotify.ScopePlaylistModifyPrivate,
			spotify.ScopePlaylistModifyPublic,
			spotify.ScopePlaylistReadPrivate,
		),
		genius.New(
			os.Getenv("GENIUS_CLIENT_ID"),
			os.Getenv("GENIUS_CLIENT_SECRET"),
			os.Getenv("GENIUS_REDIRECT_URL"),
			genius.ScopeMe,
		),
	)

	// setup server
	e := echo.New()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()

	go startServer(e, db, httpClient, store, minioClient)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func startServer(e *echo.Echo, db *sqlx.DB, httpClient *http.Client, store sessions.Store, minioClient *minio.Client) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowCredentials: true,
	}))

	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("%s, World!", os.Getenv("HELLO")))
	})

	setupAPIRouter(e, db, httpClient, store, minioClient)

	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
}

func setupAPIRouter(e *echo.Echo, db *sqlx.DB, httpClient *http.Client, store sessions.Store, minioClient *minio.Client) {
	apiRouter := e.Group("/api")

	apiRouter.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "You have been authenticated")
	})
	playlistRouter := apiRouter.Group("/playlists")
	searchRouter := apiRouter.Group("/search")
	oauthRouter := apiRouter.Group("/oauth")
	metadataRouter := apiRouter.Group("/metadata")

	setupPlaylistRoutes(playlistRouter, db, store, minioClient)
	setupSearchRoutes(searchRouter, httpClient)
	setupOAuthRoutes(oauthRouter, store)
	setupMetadataRoutes(metadataRouter)
}

func setupPlaylistRoutes(router *echo.Group, db *sqlx.DB, store sessions.Store, minioClient *minio.Client) {
	// setup playlist endpoint
	playlistRepository := repository.NewPlaylistRepository(db, minioClient)
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
	router.POST("/:id/image", playlistHandler.UploadPictureForPlaylist)

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
	router.GET("/check_auth/:provider", oauthHandler.CheckAuthHandler)
	router.GET("/logout/:provider", oauthHandler.LogoutHandler)
}

func setupMetadataRoutes(router *echo.Group) {
	metadataService := service.NewMetadataService()
	metadataHandler := rest.NewMetadataHandler(metadataService)

	router.POST("/song_lyrics", metadataHandler.GetLyrics)
}
