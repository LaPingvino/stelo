// Package sound plays the game's sounds through Ebitengine's audio: the
// original's effects and music, converted for today's machines by
// scripts/convert-audio.sh into Ogg Opus at SampleRate, decoded by gopus (a
// pure-Go Opus decoder) (audio/effects,
// audio/music; the 2006 WAV and MIDI files stay as they were).
//
// The game names a sound as action_sound does: an effect by its name; music
// as "@music:name" (sound_loop) or "@sound:name" (played once). As in Game
// Maker 5, where MIDI went through one music channel, starting a tune stops
// the one before; effects may sound together.
package sound

import (
	"bytes"
	"io"
	"io/fs"
	"path"
	"strings"
	"sync"

	"github.com/LaPingvino/gopus"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// SampleRate is the audio context's rate, and the converted files' rate, so
// nothing is resampled at run time.
const SampleRate = 48000 // Opus' own rate

// Dir is where the converted sounds are (effects/, music/).
const Dir = "audio"

// Decode reads an effect into PCM at SampleRate.
func Decode(fsys fs.FS, name string) ([]byte, error) {
	data, err := fs.ReadFile(fsys, path.Join(Dir, "effects", name+".opus"))
	if err != nil {
		return nil, err
	}
	s, err := gopus.NewStream(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Player plays the effects, each decoded once — all of them in the
// background as soon as the player is made, so no first play has to decode
// inside a game step — and one tune at a time, streamed as it plays.
type Player struct {
	fsys fs.FS
	ctx  *audio.Context
	mu   sync.Mutex
	pcm  map[string][]byte
	done chan struct{} // closed when the background preload has finished

	music     *audio.Player
	musicName string
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
		names, _ := fs.Glob(p.fsys, Dir+"/effects/*.opus")
		for _, n := range names {
			p.decoded(strings.TrimSuffix(path.Base(n), ".opus"))
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

// Play is action_sound: an effect from the start, or a tune ("@music:" looped,
// "@sound:" once) in place of the one playing.
func (p *Player) Play(name string) {
	if tune, ok := strings.CutPrefix(name, "@music:"); ok {
		p.playMusic(tune, true)
		return
	}
	if tune, ok := strings.CutPrefix(name, "@sound:"); ok {
		p.playMusic(tune, false)
		return
	}
	if pcm := p.decoded(name); pcm != nil {
		p.ctx.NewPlayerFromBytes(pcm).Play()
	}
}

// playMusic starts a tune in place of the one playing, streamed from its file.
func (p *Player) playMusic(name string, loop bool) {
	p.StopMusic()
	data, err := fs.ReadFile(p.fsys, path.Join(Dir, "music", name+".opus"))
	if err != nil {
		return
	}
	s, err := gopus.NewStream(bytes.NewReader(data))
	if err != nil {
		return
	}
	var src io.Reader = s
	if loop {
		src = audio.NewInfiniteLoop(s, s.Length())
	}
	pl, err := p.ctx.NewPlayer(src)
	if err != nil {
		return
	}
	p.music, p.musicName = pl, name
	pl.Play()
}

// StopMusic stops the tune playing (sound_stop_all for the music channel).
func (p *Player) StopMusic() {
	if p.music != nil {
		p.music.Close()
		p.music, p.musicName = nil, ""
	}
}

// Playing is sound_isplaying for a tune: whether it is the one playing.
func (p *Player) Playing(name string) bool {
	return p.music != nil && p.musicName == name && p.music.IsPlaying()
}
