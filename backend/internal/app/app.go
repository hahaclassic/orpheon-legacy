package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	audioconverter "github.com/hahaclassic/orpheon/backend/internal/adapters/audio-converter"
	bcrypt_hasher "github.com/hahaclassic/orpheon/backend/internal/adapters/password-hasher/bcrypt-hasher"
	jwttokens "github.com/hahaclassic/orpheon/backend/internal/adapters/tokens/jwt"
	"github.com/hahaclassic/orpheon/backend/internal/config"
	auth_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/auth"
	album_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/album"
	artist_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/artist"
	genre_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/genre"
	license_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/license"
	playlist_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/playlist"
	search_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/search"
	track_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/content/track"
	stats_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/stat"
	user_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/http/api/user"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/middleware"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/router"
	album_router "github.com/hahaclassic/orpheon/backend/internal/controller/http/router/router-registrators/album"
	artist_router "github.com/hahaclassic/orpheon/backend/internal/controller/http/router/router-registrators/artist"
	playlist_router "github.com/hahaclassic/orpheon/backend/internal/controller/http/router/router-registrators/playlist"
	track_router "github.com/hahaclassic/orpheon/backend/internal/controller/http/router/router-registrators/track"
	user_me_router "github.com/hahaclassic/orpheon/backend/internal/controller/http/router/router-registrators/user-me"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/utils/cookie"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/auth"
	content_aggregator "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/aggregator"
	album_cover_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/cover"
	album_meta_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/meta"
	album_tracks_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/tracks"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/assign"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/avatar"
	artist_meta_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/meta"
	genre_assign "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/genre/assign"
	genre_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/genre/meta"
	license_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/license"
	playlist_aggregator "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/aggregator"
	playlist_cover_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/cover"
	playlist_deletion_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/deleter"
	playlist_favorites_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/favorites"
	playlist_meta_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/meta"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/policy"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/privacy"
	playlist_tracks_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/playlist/tracks"
	search_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/search"
	audio_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/track/audio"
	track_meta_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/track/meta"
	tracksegment "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/track/segment"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/stat/processor"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/user"
	minio_client "github.com/hahaclassic/orpheon/backend/internal/infrastructure/minio"
	"github.com/hahaclassic/orpheon/backend/internal/infrastructure/postgres"
	"github.com/hahaclassic/orpheon/backend/internal/infrastructure/redis"
	auth_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/auth/auth-repo/postgres"
	refresh_redis "github.com/hahaclassic/orpheon/backend/internal/repository/auth/refresh-token/redis"
	album_cover_minio "github.com/hahaclassic/orpheon/backend/internal/repository/content/album/cover/minio"
	album_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/album/meta/postgres"
	album_tracks_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/album/tracks/postgres"
	assign_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/artist/assign/postgres"
	avatar_minio "github.com/hahaclassic/orpheon/backend/internal/repository/content/artist/avatar/minio"
	artist_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/artist/meta/postgres"
	genre_assign_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/genre/assign/postgres"
	genre_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/genre/meta/postgres"
	license_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/license/postgres"
	access_cache_local "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-cache/local"
	access_cache_redis "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-cache/redis"
	access_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-meta/default/postgres"
	access_meta "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-meta/with-cache"
	playlist_cover_minio "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/cover/minio"
	favorites_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/favorites/postgres"
	playlist_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/meta/postgres"
	playlist_tracks_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/tracks/postgres"
	search_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/search/postgres"
	audio_fs "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/audio/fs"
	audio_minio "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/audio/minio"
	track_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/meta/postgres"
	segment_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/segment/postgres"
	user_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/user/postgres"
	"github.com/minio/minio-go/v7"
)

func Run(conf *config.Config) {
	ctx := context.Background()

	pgxpool := postgres.NewPostgresPool(conf.Postgres)
	defer pgxpool.Close()

	redisClient, err := redis.NewRedisClient(conf.Redis)
	if err != nil {
		slog.Error("redisClient", "err", err)
		return
	}
	defer redisClient.Close()

	minioClient, err := minio_client.NewMinioClient(conf.MinIO)
	if err != nil {
		slog.Error("failed to create minio client", "err", err)
		return
	}

	// Initialize repositories
	authRepo := auth_postgres.NewAuthRepository(pgxpool)
	userRepo := user_postgres.NewUserRepository(pgxpool)
	refreshRepo := refresh_redis.NewRefreshTokenRepository(redisClient, &conf.RefreshToken)
	trackRepo := track_meta_postgres.NewTrackMetaRepository(pgxpool)
	albumMetaRepo := album_meta_postgres.NewAlbumRepository(pgxpool)
	albumTrackRepo := album_tracks_postgres.NewAlbumTrackRepository(pgxpool)
	artistMetaRepo := artist_meta_postgres.NewArtistMetaRepository(pgxpool)
	artistAssignRepo := assign_postgres.NewArtistAssignRepository(pgxpool)
	genreRepo := genre_meta_postgres.NewGenreRepository(pgxpool)
	genreAssignRepo := genre_assign_postgres.NewGenreAssignRepository(pgxpool)
	licenseRepo := license_postgres.NewLicenseRepository(pgxpool)
	searchRepo := search_postgres.NewSearchRepository(pgxpool)
	playlistRepo := playlist_meta_postgres.NewPlaylistMetaRepository(pgxpool)
	playlistTrackRepo := playlist_tracks_postgres.NewPlaylistTracksRepository(pgxpool)
	playlistFavoriteRepo := favorites_postgres.NewPlaylistFavoriteRepository(pgxpool)
	segmentRepo := segment_postgres.NewTrackSegmentRepository(pgxpool)

	albumCoverRepo, err := album_cover_minio.NewAlbumCoverRepository(ctx, minioClient, conf.MinIO.BucketAlbum)
	if err != nil {
		slog.Error("failed to create album cover repository", "err", err)
		return
	}

	artistAvatarRepo, err := avatar_minio.NewArtistAvatarRepository(ctx, minioClient, conf.MinIO.BucketArtistAvatar)
	if err != nil {
		slog.Error("failed to create artist avatar repository", "err", err)
		return
	}

	audioRepo, err := setupAudioStorage(ctx, conf, minioClient)
	if err != nil {
		slog.Error("failed to create audio file repository", "err", err)
		return
	}

	// audioRepo, err := audio_minio.NewAudioFileRepository(ctx, minioClient, conf.MinIO.BucketAudio)
	// if err != nil {
	// 	slog.Error("failed to create audio file repository", "err", err)
	// 	return
	// }

	playlistCoverRepo, err := playlist_cover_minio.NewPlaylistCoverRepository(ctx, minioClient, conf.MinIO.BucketPlaylist)
	if err != nil {
		slog.Error("failed to create playlist cover repository", "err", err)
		return
	}

	playlistAccessRepo := access_meta_postgres.NewPlaylistAccessRepository(pgxpool)
	accessCacheLocal, err := access_cache_local.NewAccessCache(conf.LocalAccessMetaCache.Size)
	if err != nil {
		slog.Error("failed to create access cache local", "err", err)
		return
	}
	accessMetaCacheRedis := access_cache_redis.NewAccessCache(redisClient, &conf.RedisAccessMetaCache)
	playlistAccessRepoWithCache := access_meta.New(playlistAccessRepo,
		access_meta.WithL1Cache(accessCacheLocal),
		access_meta.WithL2Cache(accessMetaCacheRedis))

	// Initialize services
	hasher := bcrypt_hasher.New(conf.PasswordHasher.Cost)
	tokenService := jwttokens.New(conf.AccessToken)
	userService := user.New(userRepo)
	authService := auth.NewAuthService(authRepo, refreshRepo, userService, hasher, tokenService)
	playlistPolicyService := policy.New(playlistAccessRepoWithCache)

	// Initialize content services
	segmentService := tracksegment.NewTrackSegmentService(segmentRepo)
	trackService := track_meta_service.NewTrackMetaService(trackRepo, segmentService)
	trackAudioService := audio_service.New(audioRepo, audioconverter.New())
	artistMetaService := artist_meta_service.New(artistMetaRepo)
	playlistMetaService := playlist_meta_service.NewPlaylistMetaService(playlistRepo, playlistPolicyService, playlistAccessRepo)
	playlistTrackService := playlist_tracks_service.NewPlaylistTrackService(playlistTrackRepo, playlistPolicyService)
	playlistFavoriteService := playlist_favorites_service.NewPlaylistFavoriteService(playlistFavoriteRepo, playlistPolicyService)
	playlistCoverService := playlist_cover_service.New(playlistCoverRepo, playlistPolicyService)
	playlistDeletionService := playlist_deletion_service.New(
		playlist_deletion_service.WithMetaDeletion(playlistMetaService),
		playlist_deletion_service.WithCoverDeletion(playlistCoverService),
		playlist_deletion_service.WithTracksDeletion(playlistTrackService),
		playlist_deletion_service.WithFavoritesDeletion(playlistFavoriteService),
	)
	playlistPrivacyService := privacy.NewPlaylistPrivacyChanger(
		playlistPolicyService,
		playlistFavoriteService,
		playlistAccessRepo,
	)
	genreService := genre_service.NewGenreService(genreRepo)
	genreAssignService := genre_assign.NewGenreAssignService(genreAssignRepo)
	licenseService := license_service.NewLicenseService(licenseRepo)
	albumMetaService := album_meta_service.New(albumMetaRepo)
	albumCoverService := album_cover_service.New(albumCoverRepo)
	albumTrackService := album_tracks_service.NewAlbumTrackService(albumTrackRepo)
	artistAssignService := assign.NewArtistAssignService(artistAssignRepo)
	artistAvatarService := avatar.NewArtistCoverService(artistAvatarRepo)
	searchService := search_service.NewSearchService(searchRepo)
	listeningStatService := processor.NewListeningStatService(trackRepo, segmentRepo)

	contentAggregator := content_aggregator.NewContentAggregator(
		trackService,
		artistAssignService,
		albumMetaService,
		licenseService,
		genreService,
	)

	playlistAggregator := playlist_aggregator.NewPlaylistAggregator(
		playlistMetaService,
		playlistFavoriteService,
		playlistTrackService,
		userService,
	)

	cookieTokensSetter := cookie.NewCookieTokensSetter(&conf.Cookie)

	authMiddleware := middleware.NewAuthMiddleware(authService, cookieTokensSetter)
	authMiddlewareRequired := authMiddleware.Optional() //authMiddleware.Required()
	authMiddlewareOptional := authMiddleware.Optional()

	authController := auth_ctrl.NewAuthController(authService, cookieTokensSetter, authMiddlewareRequired)

	genreController := genre_ctrl.NewGenreController(genreService, authMiddlewareRequired)
	genreAssignController := genre_ctrl.NewGenreAssignController(genreAssignService)
	licenseController := license_ctrl.NewLicenseController(licenseService, authMiddlewareRequired)
	artistMetaController := artist_ctrl.NewArtistMetaController(artistMetaService)
	artistAvatarController := artist_ctrl.NewArtistAvatarController(artistAvatarService)
	artistAssignController := artist_ctrl.NewArtistAssignController(artistAssignService, contentAggregator)
	albumMetaController := album_ctrl.NewAlbumMetaController(albumMetaService, contentAggregator)
	albumTrackController := album_ctrl.NewAlbumTrackController(albumTrackService, contentAggregator)
	albumCoverController := album_ctrl.NewAlbumCoverController(albumCoverService)
	trackMetaController := track_ctrl.NewTrackMetaController(trackService, contentAggregator)
	trackAudioController := track_ctrl.NewTrackAudioController(trackAudioService)
	searchController := search_ctrl.NewSearchController(searchService, contentAggregator, playlistAggregator, authMiddlewareOptional)
	userController := user_ctrl.NewUserController(userService)
	playlistMetaController := playlist_ctrl.NewPlaylistMetaController(playlistMetaService,
		playlistDeletionService,
		playlistPrivacyService,
		playlistAggregator,
	)
	playlistCoverController := playlist_ctrl.NewPlaylistCoverController(playlistCoverService)
	playlistTrackController := playlist_ctrl.NewPlaylistTrackController(playlistTrackService, contentAggregator)
	playlistFavoriteController := playlist_ctrl.NewPlaylistFavoritesController(playlistFavoriteService, playlistAggregator)
	trackSegmentController := track_ctrl.NewTrackSegmentController(segmentService)
	statController := stats_ctrl.NewStatController(listeningStatService)

	albumRouter := album_router.NewAlbumRouter(
		albumMetaController, albumCoverController,
		albumTrackController, genreAssignController, authMiddlewareRequired)

	artistRouter := artist_router.NewArtistRouter(
		artistMetaController, artistAvatarController,
		artistAssignController, authMiddlewareRequired)

	playlistRouter := playlist_router.NewPlaylistRouter(
		playlistMetaController, playlistTrackController, playlistCoverController, authMiddlewareRequired)

	trackRouter := track_router.NewTrackRouter(trackMetaController,
		trackSegmentController, trackAudioController, statController, artistAssignController, authMiddlewareRequired)

	meRouter := user_me_router.NewMeRouter(playlistMetaController, userController,
		playlistFavoriteController, authMiddlewareRequired)

	loggerMiddleware, err := middleware.SetupLoggerMiddleware(conf.Logger.Path, conf.Logger.Level)
	if err != nil {
		slog.Error("failed to create logger middleware", "err", err)
		return
	}

	// Initialize router
	router := router.SetupRouter(
		[]router.RoutersRegistrator{
			authController,
			genreController,
			licenseController,
			searchController,

			albumRouter,
			artistRouter,
			playlistRouter,
			trackRouter,
			meRouter,
		},
		[]gin.HandlerFunc{
			loggerMiddleware,
			middleware.CORSMiddleware(),
		})

	// addr := net.JoinHostPort(conf.HTTP.Host, conf.HTTP.Port)
	// if err := router.Run(addr); err != nil {
	// 	slog.Error("failed to start HTTP server", "err", err)
	// }
	addr := net.JoinHostPort(conf.HTTP.Host, conf.HTTP.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("starting server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	} else {
		slog.Info("server exited properly")
	}
}

func setupAudioStorage(ctx context.Context, conf *config.Config, minioClient *minio.Client) (audio_service.AudioFileRepository, error) {
	var (
		err       error
		audioRepo audio_service.AudioFileRepository
	)

	switch conf.AudioStorage.Type {
	case "fs":
		audioRepo, err = audio_fs.NewAudioFileRepository(conf.AudioStorage.BasePath)
	case "minio":
		audioRepo, err = audio_minio.NewAudioFileRepository(ctx, minioClient, conf.MinIO.BucketAudio)
	default:
		return nil, fmt.Errorf("audio file repository implementation are not specified")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create audio file repository: %w", err)
	}

	return audioRepo, nil
}
