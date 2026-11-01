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

func TestLevelRoundTrip(t *testing.T) {
	dir := t.TempDir()
	Dir = func() string { return dir }
	if _, ok := LoadLevel("nivelo1"); ok {
		t.Fatal("nothing saved yet")
	}
	SaveLevel("nivelo1", "stelo-level 1\nname nivelo1\n")
	if got, ok := LoadLevel("nivelo1"); !ok || got != "stelo-level 1\nname nivelo1\n" {
		t.Errorf("loaded %q %v", got, ok)
	}
}
