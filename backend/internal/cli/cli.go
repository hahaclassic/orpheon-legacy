package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	audioconverter "github.com/hahaclassic/orpheon/backend/internal/adapters/audio-converter"
	bcrypt_hasher "github.com/hahaclassic/orpheon/backend/internal/adapters/password-hasher/bcrypt-hasher"
	jwttokens "github.com/hahaclassic/orpheon/backend/internal/adapters/tokens/jwt"
	"github.com/hahaclassic/orpheon/backend/internal/config"
	auth_cli_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/cli/api/auth"
	album_cli_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/cli/api/content/album"
	artist_cli_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/cli/api/content/artist"
	genre_cli_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/cli/api/content/genre"
	license_cli_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/cli/api/content/license"
	playlist_cli_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/cli/api/content/playlist"
	search_cli_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/cli/api/content/search"
	track_cli_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/cli/api/content/track"
	player_cli_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/cli/api/player"
	user_cli_ctrl "github.com/hahaclassic/orpheon/backend/internal/controller/cli/api/user"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/player"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/auth"
	album_cover_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/cover"
	album_meta_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/meta"
	album_tracks_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/album/tracks"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/assign"
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/avatar"
	artist_meta_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/artist/meta"
	genre_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/genre/meta"
	license_service "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/license"
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
	"github.com/hahaclassic/orpheon/backend/internal/domain/services/user"
	"github.com/hahaclassic/orpheon/backend/internal/infrastructure/minio"
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
	genre_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/genre/meta/postgres"
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
	audio_minio "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/audio/minio"
	track_meta_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/meta/postgres"
	segment_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/content/track/segment/postgres"
	user_postgres "github.com/hahaclassic/orpheon/backend/internal/repository/user/postgres"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
	tableoutput "github.com/hahaclassic/orpheon/backend/pkg/table"
)

func Run(conf *config.Config) {
	ctx := context.Background()
	_ = session.Instance()

	pgxpool := postgres.NewPostgresPool(conf.Postgres)
	defer pgxpool.Close()

	redisClient, err := redis.NewRedisClient(conf.Redis)
	if err != nil {
		slog.Error("redisClient", "err", err)
		return
	}
	defer redisClient.Close()

	minioClient, err := minio.NewMinioClient(conf.MinIO)
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
	genreRepo := genre_postgres.NewGenreRepository(pgxpool)
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

	audioRepo, err := audio_minio.NewAudioFileRepository(ctx, minioClient, conf.MinIO.BucketAudio)
	if err != nil {
		slog.Error("failed to create audio file repository", "err", err)
		return
	}

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
	licenseService := license_service.NewLicenseService(licenseRepo)
	albumMetaService := album_meta_service.New(albumMetaRepo)
	albumCoverService := album_cover_service.New(albumCoverRepo)
	albumTrackService := album_tracks_service.NewAlbumTrackService(albumTrackRepo)
	artistAssignService := assign.NewArtistAssignService(artistAssignRepo)
	artistAvatarService := avatar.NewArtistCoverService(artistAvatarRepo)
	searchService := search_service.NewSearchService(searchRepo)
	//listeningStatService := processor.NewListeningStatService(trackRepo, segmentRepo)

	authController := auth_cli_ctrl.NewAuthController(authService)
	genreController := genre_cli_ctrl.NewGenreController(genreService)
	licenseController := license_cli_ctrl.NewLicenseController(licenseService)
	artistMetaController := artist_cli_ctrl.NewArtistMetaController(artistMetaService)
	artistAvatarController := artist_cli_ctrl.NewArtistAvatarController(artistAvatarService)
	artistAssignController := artist_cli_ctrl.NewArtistAssignController(artistAssignService)
	albumMetaController := album_cli_ctrl.NewAlbumMainController(albumMetaService, albumCoverService)
	albumTrackController := album_cli_ctrl.NewAlbumTrackController(albumTrackService)
	albumCoverController := album_cli_ctrl.NewAlbumCoverController(albumCoverService)
	trackMetaController := track_cli_ctrl.NewTrackMetaController(trackService)
	trackAudioController := track_cli_ctrl.NewTrackAudioController(trackAudioService)
	searchController := search_cli_ctrl.NewSearchController(searchService)
	userController := user_cli_ctrl.NewUserController(userService)
	playlistMetaController := playlist_cli_ctrl.NewPlaylistMetaController(playlistMetaService,
		playlistPrivacyService, playlistDeletionService,
	)
	playlistCoverController := playlist_cli_ctrl.NewPlaylistCoverController(playlistCoverService)
	playlistTrackController := playlist_cli_ctrl.NewPlaylistTrackController(playlistTrackService)
	playlistFavoriteController := playlist_cli_ctrl.NewPlaylistFavoriteController(playlistFavoriteService)
	trackSegmentController := track_cli_ctrl.NewTrackSegmentController(segmentService)

	player := player.NewPlayer(trackAudioService)
	playerController := player_cli_ctrl.NewPlayerController(player, albumTrackService,
		playlistTrackService, trackService)

	tablePrinter := tableoutput.NewTablePrinter()
	router := cmdrouter.NewCmdRouter("Orpheon", tablePrinter)

	contentGroup := router.Group("Content")

	trackGroup := contentGroup.Group("Tracks")
	albumGroup := contentGroup.Group("Albums")
	artistGroup := contentGroup.Group("Artists")
	contentGroup.Group("Genres", genreController.Menu()...)
	contentGroup.Group("Licenses", licenseController.Menu()...)
	contentGroup.Group("Search", searchController.Menu()...)

	trackGroup.Group("Meta", trackMetaController.Menu()...)
	trackGroup.Group("Audio", trackAudioController.Menu()...)
	trackGroup.SetOptionHandlers(trackSegmentController.Menu()...)

	albumGroup.Group("Meta", albumMetaController.Menu()...)
	albumGroup.Group("Covers", albumCoverController.Menu()...)
	albumGroup.SetOptionHandlers(albumTrackController.Menu()...)

	artistGroup.Group("Meta", artistMetaController.Menu()...)
	artistGroup.Group("Artist's albums", artistAssignController.AlbumsMenu()...)
	artistGroup.Group("Artist's tracks", artistAssignController.TracksMenu()...)
	artistGroup.Group("Avatar", artistAvatarController.Menu()...)

	libraryGroup := router.Group("Playlists")
	libraryGroup.Group("Meta", playlistMetaController.Menu()...)
	libraryGroup.Group("Covers", playlistCoverController.Menu()...)
	libraryGroup.Group("Tracks", playlistTrackController.Menu()...)
	libraryGroup.Group("Favorites", playlistFavoriteController.Menu()...)

	router.Group("Player", playerController.Menu()...)

	router.Group("User Profile", userController.Menu()...)
	_ = router.Group("Auth Menu", authController.Menu()...)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Println("Orpheon. CLI")
	if err := authController.RefreshToken(ctx); err == nil {
		userInfo, err := userService.GetUser(ctx, session.Claims().UserID)
		if err == nil {
			fmt.Printf("You are logged in as '%s'\n", userInfo.Name)
		}
	}

	go func() {
		router.Run(ctx)
		stop()
	}()

	<-ctx.Done()
	slog.Info("Orpheon. CLI exited")
}
