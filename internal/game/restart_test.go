package game

import (
	"testing"

	"github.com/LaPingvino/stelo"
)

// The pause menu's restart: the level again, with the score and lives it
// began with.
func TestRestartLevel(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(stelo.Files, a, "nivelo1")
	if err != nil {
		t.Fatal(err)
	}
	if g.Level() != 1 {
		t.Fatalf("Level() = %d, want 1", g.Level())
	}
	score, lives := g.G.Score, g.G.Lives
	old := g.World
	g.G.Score, g.G.Lives = score+500, lives-1
	if err := g.RestartLevel(); err != nil {
		t.Fatal(err)
	}
	if g.World == old || g.Level() != 1 || g.G.Score != score || g.G.Lives != lives {
		t.Fatalf("after restart: new world %v, level %d, score %d (want %d), lives %d (want %d)",
			g.World != old, g.Level(), g.G.Score, score, g.G.Lives, lives)
	}
}
