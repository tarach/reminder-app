package audio

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type countingStreamer struct {
	s      beep.Streamer
	pos    int
	length int
}

func (c *countingStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = c.s.Stream(samples)
	c.pos += n
	if c.length > 0 && c.pos >= c.length {
		c.pos %= c.length
	}
	return n, ok
}
func (c *countingStreamer) Err() error { return nil }

// Player handles audio playback logic without UI coupling
type Player struct {
	counter      *countingStreamer
	ctrl         *beep.Ctrl
	format       beep.Format
	bufLenFrames int
	volFx        *effects.Volume
	mu           sync.RWMutex
	isPlaying    bool
	initOnce     sync.Once
}

// NewPlayer creates and initializes a new audio player
func NewPlayer(audioData []byte) (*Player, error) {
	p := &Player{}
	if err := p.initialize(audioData); err != nil {
		return nil, fmt.Errorf("failed to initialize player: %w", err)
	}
	return p, nil
}

func (p *Player) initialize(audioData []byte) error {
	rc := io.NopCloser(bytes.NewReader(audioData))
	dec, format, err := mp3.Decode(rc)
	if err != nil {
		return fmt.Errorf("failed to decode MP3: %w", err)
	}
	p.format = format

	buf := beep.NewBuffer(p.format)
	buf.Append(dec)
	_ = dec.Close()

	p.bufLenFrames = buf.Len()

	if err := speaker.Init(p.format.SampleRate, p.format.SampleRate.N(time.Second/10)); err != nil {
		return fmt.Errorf("failed to initialize speaker: %w", err)
	}

	loop := beep.Loop(-1, buf.Streamer(0, p.bufLenFrames))
	p.counter = &countingStreamer{s: loop, length: p.bufLenFrames}

	p.volFx = &effects.Volume{
		Streamer: p.counter,
		Base:     10,
		Volume:   0.0,
		Silent:   false,
	}

	p.ctrl = &beep.Ctrl{Streamer: p.volFx, Paused: true}
	speaker.Play(p.ctrl)

	return nil
}

// Play starts audio playback
func (p *Player) Play() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.isPlaying {
		p.ctrl.Paused = false
		p.isPlaying = true
	}
}

// Stop pauses audio playback
func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isPlaying {
		p.ctrl.Paused = true
		p.isPlaying = false
	}
}

// IsPlaying returns the current playback state
func (p *Player) IsPlaying() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isPlaying
}

// SetVolume sets the volume (0-100)
func (p *Player) SetVolume(volume float64) {
	speaker.Lock()
	defer speaker.Unlock()

	if volume <= 0 {
		p.volFx.Silent = true
		p.volFx.Volume = -10
		return
	}
	p.volFx.Silent = false

	// Map 0-100 to a usable dB range
	minDB := -2.0
	maxDB := 0.0

	// Linear mapping from slider (1-100) to dB range
	normalizedVolume := volume / 100.0
	p.volFx.Volume = minDB + (maxDB-minDB)*normalizedVolume
}

// GetPosition returns the current position and total duration in frames
func (p *Player) GetPosition() (pos, total int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	total = p.bufLenFrames
	if total > 0 && p.counter != nil {
		pos = p.counter.pos % total
	}
	return pos, total
}

// GetDuration converts frames to duration
func (p *Player) GetDuration(frames int) time.Duration {
	return time.Duration(float64(frames)/float64(p.format.SampleRate)) * time.Second
}

// Close cleans up resources
func (p *Player) Close() {
	p.Stop()
	speaker.Clear()
}
