package sound

import (
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
