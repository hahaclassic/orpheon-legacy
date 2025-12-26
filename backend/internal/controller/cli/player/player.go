package player

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/hahaclassic/orpheon/backend/internal/controller/cli/output"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
)

const (
	chunkSize      = 64 * 1024              // 64KB
	initialBuffer  = 10 * chunkSize         // 640KB
	bufferDuration = 100 * time.Millisecond // speaker buffer duration
)

type streamBuffer struct {
	buf  *bytes.Buffer
	done chan struct{}
}

func newStreamBuffer() *streamBuffer {
	return &streamBuffer{
		buf:  bytes.NewBuffer(nil),
		done: make(chan struct{}),
	}
}

func (s *streamBuffer) Read(p []byte) (int, error) {
	for {
		n, err := s.buf.Read(p)
		if err == io.EOF {
			select {
			case <-s.done:
				return n, io.EOF
			default:
				time.Sleep(10 * time.Millisecond)
				continue
			}
		}
		return n, err
	}
}

func (s *streamBuffer) Write(data []byte) {
	s.buf.Write(data)
}

func (s *streamBuffer) Close() error {
	select {
	case <-s.done:
	default:
		close(s.done)
	}
	return nil
}

type Player struct {
	Queue                []*entity.TrackMeta
	Current              int
	CurrentSecond        int
	IsPlaying            bool
	audioFileService     track.AudioFileService
	streamer             beep.StreamSeekCloser
	ctrl                 *beep.Ctrl
	format               beep.Format
	done                 chan struct{}
	mu                   sync.Mutex
	cancelPlayback       context.CancelFunc
	progressTickerCancel context.CancelFunc
}

func NewPlayer(audioFileService track.AudioFileService) *Player {
	return &Player{
		audioFileService: audioFileService,
		Queue:            []*entity.TrackMeta{},
	}
}

func (c *Player) Play(ctx context.Context) {
	c.mu.Lock()
	if c.IsPlaying || c.Current >= len(c.Queue) {
		c.mu.Unlock()
		return
	}
	track := c.Queue[c.Current]
	c.done = make(chan struct{})
	c.IsPlaying = true
	c.CurrentSecond = 0
	c.mu.Unlock()

	output.PrintTrack(track)
	//fmt.Println("Playing track:", track.Name, "with duration:", track.Duration)

	ctx, cancel := context.WithCancel(ctx)
	c.mu.Lock()
	c.cancelPlayback = cancel
	c.mu.Unlock()

	sb := newStreamBuffer()
	ready := make(chan struct{}, 1)

	go c.streamAudio(ctx, track, sb, ready, int64(c.CurrentSecond))
	<-ready

	if err := c.startPlayback(sb); err != nil {
		log.Printf("playback error: %v", err)
		return
	}

	go c.trackProgress(ctx, track)
}

func (c *Player) streamAudio(ctx context.Context, track *entity.TrackMeta, sb *streamBuffer, ready chan struct{}, startSecond int64) {
	var offset int64 = startSecond * 44100 * 4 // 44.1kHz * 4 bytes per frame (estimate)
	var total int

	for {
		select {
		case <-ctx.Done():
			sb.Close()
			return
		default:
		}

		chunk := &entity.AudioChunk{
			TrackID: track.ID,
			Start:   offset,
			End:     offset + chunkSize,
		}

		resp, err := c.audioFileService.GetAudioChunk(ctx, chunk)
		if err != nil || len(resp.Data) == 0 {
			break
		}

		sb.Write(resp.Data)
		offset += int64(len(resp.Data))
		total += len(resp.Data)

		if total >= initialBuffer {
			select {
			case ready <- struct{}{}:
			default:
			}
		}
	}
	sb.Close()
}

func (c *Player) startPlayback(sb *streamBuffer) error {
	streamer, format, err := mp3.Decode(sb)
	if err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(bufferDuration)); err != nil {
		return fmt.Errorf("speaker init error: %w", err)
	}

	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}

	c.mu.Lock()
	c.streamer = streamer
	c.ctrl = ctrl
	c.format = format
	c.mu.Unlock()

	speaker.Play(beep.Seq(ctrl, beep.Callback(func() {
		c.mu.Lock()
		c.IsPlaying = false
		c.mu.Unlock()
		close(c.done)
	})))

	return nil
}

func (c *Player) trackProgress(ctx context.Context, track *entity.TrackMeta) {
	tickerCtx, cancel := context.WithCancel(ctx)
	c.mu.Lock()
	c.progressTickerCancel = cancel
	c.mu.Unlock()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-tickerCtx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			if c.IsPlaying {
				c.CurrentSecond++
				if c.CurrentSecond >= int(track.Duration) {
					c.CurrentSecond = int(track.Duration)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *Player) Pause() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ctrl != nil && c.IsPlaying {
		c.ctrl.Paused = true
		c.IsPlaying = false
	}
}

func (c *Player) Resume() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ctrl != nil && !c.IsPlaying {
		fmt.Println("Resuming track:", c.Queue[c.Current].Name)
		c.ctrl.Paused = false
		c.IsPlaying = true
	}
}

func (c *Player) Stop() {
	c.mu.Lock()
	if !c.IsPlaying && c.streamer == nil {
		c.mu.Unlock()
		return
	}
	if c.cancelPlayback != nil {
		c.cancelPlayback()
	}
	if c.progressTickerCancel != nil {
		c.progressTickerCancel()
		c.progressTickerCancel = nil
	}
	if c.ctrl != nil {
		c.ctrl.Paused = true
	}
	if c.streamer != nil {
		c.streamer.Close()
		c.streamer = nil
	}
	c.ctrl = nil
	c.IsPlaying = false
	if c.done != nil {
		select {
		case <-c.done:
		default:
			close(c.done)
		}
	}
	c.mu.Unlock()
}

func (c *Player) Next() {
	c.Stop()
	c.mu.Lock()
	if c.Current < len(c.Queue)-1 {
		c.Current++
		c.CurrentSecond = 0
	}
	c.mu.Unlock()
	go c.Play(context.Background())
}

func (c *Player) Previous() {
	c.Stop()
	c.mu.Lock()
	if c.Current > 0 {
		c.Current--
		c.CurrentSecond = 0
	}
	c.mu.Unlock()
	go c.Play(context.Background())
}

func (c *Player) AddToQueue(tracks []*entity.TrackMeta) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.Queue) > 0 {
		c.Queue = c.Queue[:c.Current+1]
	}
	c.Queue = append(c.Queue, tracks...)
}

func (c *Player) SeekTo(second int) {
	c.mu.Lock()
	if c.Current >= len(c.Queue) || second < 0 || second >= int(c.Queue[c.Current].Duration) {
		c.mu.Unlock()
		return
	}
	c.CurrentSecond = second
	c.mu.Unlock()

	c.Stop()
	go c.Play(context.Background())
}
