// Package sound plays the original's sound effects (original/sounds/*.wav)
// through Ebitengine's audio. The game names a sound as action_sound does;
// "@music:…" and "@sound:…" (the MIDI tunes) are left to the music player.
// PARTIAL: music — the 2006 MIDI files are to be rendered to OGG.
package sound

import (
	"bytes"
	"io/fs"
	"path"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// SampleRate is the audio context's rate; the WAVs are resampled to it.
const SampleRate = 44100

// Decode reads a WAV into PCM at SampleRate (the original's files, any rate).
func Decode(fsys fs.FS, name string) ([]byte, error) {
	data, err := fs.ReadFile(fsys, path.Join("original/sounds", name+".wav"))
	if err != nil {
		return nil, err
	}
	s, err := wav.DecodeWithSampleRate(SampleRate, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Player plays the effects, each decoded once. All of them are decoded in the
// background as soon as the player is made: decoding (and resampling to
// SampleRate) inside a game step made the first play of each effect stall the
// game for a moment — up to 0.4 s natively for the longest, several times
// that in the browser. The pauses the original asks for (action_sleep after
// some sounds) are separate and kept.
type Player struct {
	fsys fs.FS
	ctx  *audio.Context
	mu   sync.Mutex
	pcm  map[string][]byte
	done chan struct{} // closed when the background preload has finished
}

// NewPlayer makes the audio context (once per program) and starts decoding
// every effect in the background.
func NewPlayer(fsys fs.FS) *Player {
	p := &Player{fsys: fsys, ctx: audio.NewContext(SampleRate)}
	p.startPreload()
	return p
}

func (p *Player) startPreload() {
	p.pcm = map[string][]byte{}
	p.done = make(chan struct{})
	go func() {
		defer close(p.done)
		names, _ := fs.Glob(p.fsys, "original/sounds/*.wav")
		for _, n := range names {
			p.decoded(strings.TrimSuffix(path.Base(n), ".wav"))
		}
	}()
}

// decoded returns an effect's PCM, decoding it if no one has yet (nil when it
// cannot be read).
func (p *Player) decoded(name string) []byte {
	p.mu.Lock()
	pcm, ok := p.pcm[name]
	p.mu.Unlock()
	if ok {
		return pcm
	}
	pcm, err := Decode(p.fsys, name)
	if err != nil {
		pcm = nil
	}
	p.mu.Lock()
	p.pcm[name] = pcm
	p.mu.Unlock()
	return pcm
}

// Play is action_sound(name, loop=0): the effect from the start; several may
// sound at once, as in Game Maker.
func (p *Player) Play(name string) {
	if strings.HasPrefix(name, "@") {
		return // music (MIDI) — not yet
	}
	if pcm := p.decoded(name); pcm != nil {
		p.ctx.NewPlayerFromBytes(pcm).Play()
	}
}
