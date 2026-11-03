package game

import (
	"image"
	"strings"
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

// The level-done message: in the exit's pause, not in a death's pause.
func TestLevelDone(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(stelo.Files, a, "nivelo1")
	if err != nil {
		t.Fatal(err)
	}
	w := g.World
	p := w.Room.All("character")[0]
	ex := w.Room.All("levelexit")[0]
	p.X, p.Y = ex.X, ex.Y
	for i := 0; i < 3 && w.Room.Paused == nil; i++ {
		if err := g.Step(held{}); err != nil {
			t.Fatal(err)
		}
	}
	if !w.LevelDone() {
		t.Fatalf("paused %v, goto %q", w.Room.Paused != nil, w.Room.Goto)
	}
	w.Overlay = &Overlay{}
	w.drawHUD(image.NewRGBA(image.Rect(0, 0, ScreenW, ScreenH)), 0, 3)
	found := false
	for _, tx := range w.Overlay.Texts {
		found = found || tx.S == Tr("Nivelo finita!")
	}
	if !found {
		t.Error("no 'Nivelo finita!' in the overlay")
	}
}

// The new highscore screen: the name field, then the table.
func TestHighscoresUI(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(stelo.Files, a, "nivelo1")
	if err != nil {
		t.Fatal(err)
	}
	g.World.Overlay = &Overlay{}
	g.hs, g.G.KeyboardString = hsName, "Joop"
	g.drawHighscores(image.NewRGBA(image.Rect(0, 0, ScreenW, ScreenH)))
	var texts []string
	for _, tx := range g.World.Overlay.Texts {
		texts = append(texts, tx.S)
	}
	if got := strings.Join(texts, "|"); !strings.Contains(got, "plej bonaj ludantoj!") || !strings.Contains(got, "|Joop") {
		t.Errorf("name screen: %s", got)
	}
	g.World.Overlay = &Overlay{}
	g.Highscores.Add("Joop", 900)
	g.hs = hsTable
	g.drawHighscores(image.NewRGBA(image.Rect(0, 0, ScreenW, ScreenH)))
	if tx := g.World.Overlay.Texts; tx[1].S != "1." || tx[2].S != "Joop" || tx[3].S != "900" {
		t.Errorf("first row: %q %q %q", tx[1].S, tx[2].S, tx[3].S)
	}
}
