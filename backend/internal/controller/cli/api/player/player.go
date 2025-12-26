package player_cli_ctrl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/player"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/session"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/album"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/playlist"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
	"github.com/hahaclassic/orpheon/backend/pkg/cmdrouter"
)

type PlayerController struct {
	player         *player.Player
	albumTracks    album.AlbumTrackService
	playlistTracks playlist.PlaylistTrackService
	trackService   track.TrackMetaService
}

func NewPlayerController(player *player.Player, albumTracks album.AlbumTrackService,
	playlistTracks playlist.PlaylistTrackService, trackService track.TrackMetaService) *PlayerController {
	return &PlayerController{
		player:         player,
		albumTracks:    albumTracks,
		playlistTracks: playlistTracks,
		trackService:   trackService,
	}
}

func (c *PlayerController) Menu() []cmdrouter.OptionHandler {
	return []cmdrouter.OptionHandler{
		{
			Name: "Play",
			Run:  c.play,
		},
		{
			Name: "Pause",
			Run:  c.pause,
		},
		{
			Name: "Seek",
			Run:  c.seekTo,
		},
		{
			Name: "Next",
			Run:  c.next,
		},
		{
			Name: "Previous",
			Run:  c.previous,
		},
		{
			Name: "Listen Album",
			Run:  c.listenAlbum,
		},
		{
			Name: "Listen Playlist",
			Run:  c.listenPlaylist,
		},
		{
			Name: "Listen Track",
			Run:  c.listenTrack,
		},
	}
}

func (c *PlayerController) listenAlbum(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter album ID: ")
	scanner.Scan()
	id := scanner.Text()

	albumID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse album ID: %w", err)
	}

	tracks, err := c.albumTracks.GetAllTracks(ctx, albumID)
	if err != nil {
		return fmt.Errorf("failed to get album: %w", err)
	}

	c.player.AddToQueue(tracks)
	c.player.Play(ctx)
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (c *PlayerController) listenPlaylist(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter playlist ID: ")
	scanner.Scan()
	id := scanner.Text()

	playlistID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse playlist ID: %w", err)
	}

	tracks, err := c.playlistTracks.GetAllTracks(ctx, session.Claims(), playlistID)
	if err != nil {
		return fmt.Errorf("failed to get playlist: %w", err)
	}

	c.player.AddToQueue(tracks)
	c.player.Play(ctx)
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (c *PlayerController) listenTrack(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter track ID: ")
	scanner.Scan()
	id := scanner.Text()

	trackID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("failed to parse track ID: %w", err)
	}

	track, err := c.trackService.GetTrackMeta(ctx, trackID)
	if err != nil {
		return fmt.Errorf("failed to get track: %w", err)
	}

	c.player.AddToQueue([]*entity.TrackMeta{track})
	c.player.Play(ctx)
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (c *PlayerController) play(ctx context.Context) error {
	c.player.Resume()
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (c *PlayerController) pause(ctx context.Context) error {
	c.player.Pause()
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (c *PlayerController) next(ctx context.Context) error {
	c.player.Next()
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (c *PlayerController) previous(ctx context.Context) error {
	c.player.Previous()
	time.Sleep(5 * time.Millisecond)
	return nil
}

func (c *PlayerController) seekTo(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter position to seek: ")
	scanner.Scan()
	input := scanner.Text()

	seconds := 0 // default value
	if input != "" {
		var err error
		seconds, err = strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("invalid input: %w", err)
		}
	}

	c.player.SeekTo(seconds)
	time.Sleep(5 * time.Millisecond)
	return nil
}
