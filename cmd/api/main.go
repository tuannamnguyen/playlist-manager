package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/go-playground/validator"
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
	"gopkg.in/boj/redistore.v1"
)

type CustomValidator struct {
	validator *validator.Validate
}

type SecretsManagerValues struct {
	Host string `json:"db_host"`
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	httpClient := newHTTPClient()

	isProd, err := strconv.ParseBool(os.Getenv("IS_PROD"))
	if err != nil {
		return fmt.Errorf("parsing bool: %w", err)
	}
	log.Printf("IS_PROD value: %t\n", isProd)

	awsConfig, err := newAWSConfig(ctx)
	if err != nil {
		return err
	}

	s3Client, s3PresignClient := newS3Clients(awsConfig)

	secrets, err := getSecrets(ctx, awsConfig, "db-host")
	if err != nil {
		return err
	}

	db, err := newDB(isProd, secrets, awsConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	store, err := newSessionStore(isProd)
	if err != nil {
		return err
	}
	defer store.Close()

	setupOAuth(store)

	e := echo.New()

	if err := startGracefulServer(e, db, httpClient, store, s3Client, s3PresignClient); err != nil {
		e.Logger.Fatal(err)
	}

	return nil
}

func newHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100

	return &http.Client{
		Timeout:   time.Minute,
		Transport: transport,
	}
}

func newAWSConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to config AWS: %w", err)
	}
	return cfg, nil
}

func newS3Clients(cfg aws.Config) (*s3.Client, *s3.PresignClient) {
	s3Client := s3.NewFromConfig(cfg)
	s3PresignClient := s3.NewPresignClient(s3Client)
	return s3Client, s3PresignClient
}

func getSecrets(ctx context.Context, cfg aws.Config, secretName string) (SecretsManagerValues, error) {
	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(ctx, input)
	if err != nil {
		return SecretsManagerValues{}, fmt.Errorf("get secret failed: %w", err)
	}

	var secrets SecretsManagerValues
	if err := json.Unmarshal([]byte(*result.SecretString), &secrets); err != nil {
		return SecretsManagerValues{}, fmt.Errorf("unmarshal secret: %w", err)
	}

	return secrets, nil
}

func newDB(isProd bool, secrets SecretsManagerValues, awsConfig aws.Config) (*sqlx.DB, error) {
	var host, password string
	var err error

	if isProd {
		host = secrets.Host

		password, err = auth.BuildAuthToken(
			context.TODO(),
			fmt.Sprintf("%s:%s", host, "5432"),
			os.Getenv("AWS_DEFAULT_REGION"),
			os.Getenv("POSTGRES_USER"),
			awsConfig.Credentials,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create authentication token: %w", err)
		}
	} else {
		password = os.Getenv("POSTGRES_PASSWORD")
		host = os.Getenv("POSTGRES_HOST")
	}

	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		host,
		os.Getenv("POSTGRES_USER"),
		password,
		os.Getenv("POSTGRES_DBNAME"),
	)

	db, err := sqlx.Connect("pgx", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	log.Println("connected to postgres successfully")
	return db, nil
}

func newSessionStore(isProd bool) (*redistore.RediStore, error) {
	redisInfo := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	key := os.Getenv("SESSION_SECRET")

	store, err := redistore.NewRediStore(
		10,
		"tcp",
		redisInfo,
		os.Getenv("REDIS_PASSWORD"),
		[]byte(key),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}

	store.SetMaxAge(3600)
	store.Options.Secure = isProd

	if isProd {
		store.Options.SameSite = http.SameSiteNoneMode
	}

	return store, nil
}

func setupOAuth(store sessions.Store) {
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
			spotify.ScopeStreaming,
		),
	)
}

func startGracefulServer(
	e *echo.Echo,
	db *sqlx.DB,
	httpClient *http.Client,
	store *redistore.RediStore,
	s3Client *s3.Client,
	s3PresignClient *s3.PresignClient,
) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stop()

	go startServer(e, db, httpClient, store, s3Client, s3PresignClient)

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return e.Shutdown(shutdownCtx)
}

func startServer(e *echo.Echo, db *sqlx.DB, httpClient *http.Client, store sessions.Store, s3Client *s3.Client, s3PresignedClient *s3.PresignClient) {
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	log.Printf("frontend url: %s", os.Getenv("FRONTEND_URL"))

	e.Validator = &CustomValidator{validator: validator.New()}

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "healthcheck ok")
	})

	setupAPIRouter(e, db, httpClient, store, s3Client, s3PresignedClient)

	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		// if error here, check if there are any other apps running on the same port
		e.Logger.Fatal("shutting down the server")
	}
}

func setupAPIRouter(e *echo.Echo, db *sqlx.DB, httpClient *http.Client, store sessions.Store, s3Client *s3.Client, s3PresignedClient *s3.PresignClient) {
	apiRouter := e.Group("/api")

	apiRouter.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "You have been authenticated")
	})
	playlistRouter := apiRouter.Group("/playlists")
	searchRouter := apiRouter.Group("/search")
	oauthRouter := apiRouter.Group("/oauth")
	metadataRouter := apiRouter.Group("/metadata")

	setupPlaylistRoutes(playlistRouter, db, store, s3Client, s3PresignedClient)
	setupSearchRoutes(searchRouter, httpClient)
	setupOAuthRoutes(oauthRouter, store)
	setupMetadataRoutes(metadataRouter, store)
}

func setupPlaylistRoutes(router *echo.Group, db *sqlx.DB, store sessions.Store, s3Client *s3.Client, s3PresignedClient *s3.PresignClient) {
	// setup playlist endpoint
	playlistRepository := repository.NewPlaylistRepository(db, s3Client, s3PresignedClient)
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

	// csv endpoints
	router.GET("/:playlist_id/songs/csv", playlistHandler.GetAllSongsFromPlaylistToCsv)
	router.POST("/:playlist_id/songs/csv", playlistHandler.AddSongsToPlaylistFromCsv)
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
	router.GET("/token/:provider", oauthHandler.GetAccessTokenHandler)
	router.GET("/check_auth/:provider", oauthHandler.CheckAuthHandler)
	router.GET("/logout/:provider", oauthHandler.LogoutHandler)
}

func setupMetadataRoutes(router *echo.Group, store sessions.Store) {
	metadataService := service.NewMetadataService()
	metadataHandler := rest.NewMetadataHandler(metadataService, store)

	router.POST("/song_lyrics", metadataHandler.GetLyrics)
	router.GET("/artist_information", metadataHandler.GetArtistInformation)
}
