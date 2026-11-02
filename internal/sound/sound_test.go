package sound

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/LaPingvino/gopus"
	"github.com/LaPingvino/stelo"
)

// Every effect the game plays decodes.
func TestEffectsDecode(t *testing.T) {
	for _, n := range []string{"getmushroom", "killcharacter", "killmonster", "finishlevel"} {
		pcm, err := Decode(stelo.Files, n)
		if err != nil || len(pcm) == 0 {
			t.Errorf("%s: %d bytes, %v", n, len(pcm), err)
		}
	}
}

// The background preload decodes every effect, so no first play has to.
func TestPreloadDecodesEveryEffect(t *testing.T) {
	fsys := stelo.Files
	p := &Player{fsys: fsys}
	p.startPreload()
	<-p.done
	names, _ := fs.Glob(fsys, "audio/effects/*.opus")
	if len(names) == 0 {
		t.Fatal("no effects")
	}
	for _, n := range names {
		k := strings.TrimSuffix(path.Base(n), ".opus")
		if len(p.pcm[k]) == 0 {
			t.Errorf("%s not decoded", k)
		}
	}
	t.Logf("%d effects decoded in the background", len(names))
}

// Every 2006 sound has its converted file (scripts/convert-audio.sh): each
// WAV an effect, each MIDI a tune; and every tune decodes.
func TestEveryOriginalSoundConverted(t *testing.T) {
	origs, _ := fs.Glob(os.DirFS("../.."), "original/sounds/*")
	if len(origs) == 0 {
		t.Skip("no originals")
	}
	for _, o := range origs {
		ext := path.Ext(o)
		n := strings.TrimSuffix(path.Base(o), ext)
		var want string
		switch ext {
		case ".wav":
			want = "audio/effects/" + n + ".opus"
		case ".mid":
			want = "audio/music/" + n + ".opus"
		default:
			continue
		}
		data, err := fs.ReadFile(stelo.Files, want)
		if err != nil {
			t.Errorf("%s: no %s (run scripts/convert-audio.sh)", o, want)
			continue
		}
		if ext == ".mid" {
			s, err := gopus.NewStream(bytes.NewReader(data))
			if err != nil || s.Length() == 0 {
				t.Errorf("%s: %v", want, err)
			}
		}
	}
}

// Muting: effects are skipped, the setting reads back.
func TestMute(t *testing.T) {
	p := &Player{fsys: stelo.Files}
	p.SetMuted(true)
	if !p.Muted() || p.volume() != 0 {
		t.Fatal("not muted")
	}
	p.Play("getmushroom") // no audio context: would panic if it tried to play
	p.SetMuted(false)
	if p.Muted() || p.volume() != 1 {
		t.Fatal("still muted")
	}
}
