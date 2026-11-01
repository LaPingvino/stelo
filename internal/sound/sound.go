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

// Player plays the effects, each decoded once.
type Player struct {
	fsys fs.FS
	ctx  *audio.Context
	pcm  map[string][]byte
}

// NewPlayer makes the audio context (once per program).
func NewPlayer(fsys fs.FS) *Player {
	return &Player{fsys: fsys, ctx: audio.NewContext(SampleRate), pcm: map[string][]byte{}}
}

// Play is action_sound(name, loop=0): the effect from the start; several may
// sound at once, as in Game Maker.
func (p *Player) Play(name string) {
	if strings.HasPrefix(name, "@") {
		return // music (MIDI) — not yet
	}
	pcm, ok := p.pcm[name]
	if !ok {
		var err error
		if pcm, err = Decode(p.fsys, name); err != nil {
			p.pcm[name] = nil
			return
		}
		p.pcm[name] = pcm
	}
	if pcm == nil {
		return
	}
	p.ctx.NewPlayerFromBytes(pcm).Play()
}
