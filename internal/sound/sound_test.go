package sound

import (
	"io/fs"
	"path"
	"strings"
	"testing"

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
	names, _ := fs.Glob(fsys, "original/sounds/*.wav")
	if len(names) == 0 {
		t.Fatal("no effects")
	}
	for _, n := range names {
		k := strings.TrimSuffix(path.Base(n), ".wav")
		if len(p.pcm[k]) == 0 {
			t.Errorf("%s not decoded", k)
		}
	}
	t.Logf("%d effects decoded in the background", len(names))
}
