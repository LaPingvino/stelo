//go:build !js

package store

import (
	"testing"

	"github.com/LaPingvino/stelo/internal/game"
)

func TestHighscoresRoundTrip(t *testing.T) {
	dir := t.TempDir()
	Dir = func() string { return dir }
	if h := LoadHighscores(); h != game.ClearedHighscores() {
		t.Fatalf("nothing saved yet: %+v", h)
	}
	h := game.ClearedHighscores()
	h.Add("Joop", 1230)
	SaveHighscores(h)
	if got := LoadHighscores(); got != h {
		t.Errorf("loaded %+v, want %+v", got, h)
	}
}
