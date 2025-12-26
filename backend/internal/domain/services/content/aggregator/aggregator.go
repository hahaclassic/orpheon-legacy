package content_aggregator

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/aggregator"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/artist"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/genre"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/license"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type ContentAggregator struct {
	trackService   track.TrackMetaService
	artistService  artist.ArtistAssignService
	albumService   album.AlbumMetaService
	licenseService license.LicenseService
	genreService   genre.GenreService
}

func NewContentAggregator(trackService track.TrackMetaService,
	artistService artist.ArtistAssignService,
	albumService album.AlbumMetaService,
	licenseService license.LicenseService,
	genreService genre.GenreService) *ContentAggregator {
	return &ContentAggregator{
		trackService:   trackService,
		artistService:  artistService,
		albumService:   albumService,
		licenseService: licenseService,
		genreService:   genreService,
	}
}

func (a *ContentAggregator) GetTracksByIDs(ctx context.Context, trackIDs ...uuid.UUID) (_ []*entity.TrackMetaAggregated, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetTracksAggregatedByIDs, err)
	}()

	tracks := make([]*entity.TrackMetaAggregated, len(trackIDs))

	for i, id := range trackIDs {
		trackMeta, err := a.trackService.GetTrackMeta(ctx, id)
		if err != nil {
			return nil, err
		}

		artists, err := a.artistService.GetArtistByTrack(ctx, id)
		if err != nil {
			return nil, err
		}

		album, err := a.albumService.GetAlbum(ctx, trackMeta.AlbumID)
		if err != nil {
			return nil, err
		}

		license, err := a.licenseService.GetLicenseByID(ctx, trackMeta.LicenseID)
		if err != nil {
			return nil, err
		}

		genre, err := a.genreService.GetGenreByID(ctx, trackMeta.GenreID)
		if err != nil {
			return nil, err
		}

		tracks[i] = &entity.TrackMetaAggregated{
			ID:           trackMeta.ID,
			Name:         trackMeta.Name,
			Duration:     trackMeta.Duration,
			Explicit:     trackMeta.Explicit,
			TrackNumber:  trackMeta.TrackNumber,
			TotalStreams: trackMeta.TotalStreams,
			License:      license,
			Album:        album,
			Artists:      artists,
			Genre:        genre,
		}
	}

	return tracks, nil
}

func (a *ContentAggregator) GetAlbumsByIDs(ctx context.Context, albumIDs ...uuid.UUID) (_ []*entity.AlbumMetaAggregated, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAlbumsAggregatedByIDs, err)
	}()

	albums := make([]*entity.AlbumMetaAggregated, len(albumIDs))

	for i, id := range albumIDs {
		albumMeta, err := a.albumService.GetAlbum(ctx, id)
		if err != nil {
			return nil, err
		}

		artists, err := a.artistService.GetArtistByAlbum(ctx, id)
		if err != nil {
			return nil, err
		}

		genres, err := a.genreService.GetGenreByAlbum(ctx, id)
		if err != nil {
			return nil, err
		}

		license, err := a.licenseService.GetLicenseByID(ctx, albumMeta.LicenseID)
		if err != nil {
			return nil, err
		}

		albums[i] = &entity.AlbumMetaAggregated{
			ID:          albumMeta.ID,
			Title:       albumMeta.Title,
			Label:       albumMeta.Label,
			ReleaseDate: albumMeta.ReleaseDate,
			Artists:     artists,
			Genres:      genres,
			License:     license,
		}
	}

	return albums, nil
}

func (a *ContentAggregator) GetAlbums(ctx context.Context, albums ...*entity.AlbumMeta) (_ []*entity.AlbumMetaAggregated, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAlbumsAggregated, err)
	}()

	aggregated := make([]*entity.AlbumMetaAggregated, len(albums))

	for i, album := range albums {
		artists, err := a.artistService.GetArtistByAlbum(ctx, album.ID)
		if err != nil {
			slog.Error("aggregator.GetAlbums: artists not found for album")
		}

		genres, err := a.genreService.GetGenreByAlbum(ctx, album.ID)
		if err != nil {
			slog.Error("aggregator.GetAlbums: genres not found for album")
		}

		license, err := a.licenseService.GetLicenseByID(ctx, album.LicenseID)
		if err != nil {
			slog.Error("aggregator.GetAlbums: licenses not found for album")
		}

		aggregated[i] = &entity.AlbumMetaAggregated{
			ID:          album.ID,
			Title:       album.Title,
			Label:       album.Label,
			ReleaseDate: album.ReleaseDate,
			Artists:     artists,
			Genres:      genres,
			License:     license,
		}
	}

	return aggregated, nil
}

func (a *ContentAggregator) GetTracks(ctx context.Context, tracks ...*entity.TrackMeta) (_ []*entity.TrackMetaAggregated, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetTracksAggregated, err)
	}()

	aggregated := make([]*entity.TrackMetaAggregated, len(tracks))

	for i, track := range tracks {
		artists, err := a.artistService.GetArtistByTrack(ctx, track.ID)
		if err != nil {
			return nil, err
		}

		album, err := a.albumService.GetAlbum(ctx, track.AlbumID)
		if err != nil {
			return nil, err
		}

		license, err := a.licenseService.GetLicenseByID(ctx, track.LicenseID)
		if err != nil {
			return nil, err
		}

		genre, err := a.genreService.GetGenreByID(ctx, track.GenreID)
		if err != nil {
			return nil, err
		}

		aggregated[i] = &entity.TrackMetaAggregated{
			ID:           track.ID,
			Name:         track.Name,
			Duration:     track.Duration,
			Explicit:     track.Explicit,
			TrackNumber:  track.TrackNumber,
			TotalStreams: track.TotalStreams,
			License:      license,
			Album:        album,
			Artists:      artists,
			Genre:        genre,
		}
	}

	return aggregated, nil
}
