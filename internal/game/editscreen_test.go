package game

import (
	"testing"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/editor"
	"github.com/LaPingvino/stelo/internal/level"
)

// Picking a tool on the strip and dragging across the level paints each cell
// once; right-dragging erases; Undo takes a stroke's cells back one by one.
func TestEditScreenPaint(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := stelo.Files.ReadFile("levels/nivelo1.txt")
	l, err := level.Parse(string(data))
	if err != nil {
		t.Fatal(err)
	}
	e := NewEditScreen(a, editor.New(l))
	coin := 0
	for i, tool := range Palette {
		if tool == "stelo_poentoj" {
			coin = i
		}
	}
	e.Press(coin*editor.Cell+5, paletteY+5, false)
	if e.Tool != "stelo_poentoj" {
		t.Fatalf("tool %q", e.Tool)
	}
	// Row 2 (y 64) is open air from cell 5 in level 1: drag over cells 6..8.
	before := len(l.Instances)
	e.Press(6*32+3, 2*32+3, false)
	for x := 6 * 32; x < 9*32; x += 4 {
		e.Drag(x, 2*32+10, false)
	}
	e.Release()
	if n := len(l.Instances) - before; n != 3 {
		t.Fatalf("%d coins painted, want 3", n)
	}
	e.Press(7*32+3, 2*32+3, true) // right button: erase
	e.Release()
	if o, _ := e.Doc.At(7, 2); o != "" {
		t.Fatalf("erase left %q", o)
	}
	e.Undo()
	if o, _ := e.Doc.At(7, 2); o != "stelo_poentoj" {
		t.Fatalf("undo: %q", o)
	}
}

// Play runs the edited level: a coin painted in the editor is in the room.
func TestPlayEditedLevel(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := stelo.Files.ReadFile("levels/nivelo1.txt")
	l, _ := level.Parse(string(data))
	d := editor.New(l)
	before := 0
	for _, in := range l.Instances {
		if in.Object == "stelo_poentoj" {
			before++
		}
	}
	d.Paint(7, 2, "stelo_poentoj")
	g, err := NewGameWithLevel(stelo.Files, a, d.L)
	if err != nil {
		t.Fatal(err)
	}
	if g.Title || g.World == nil || g.World.Level.Name != "nivelo1" {
		t.Fatalf("not playing the level (title %v)", g.Title)
	}
	if n := len(g.World.Room.All("stelo_poentoj")); n != before+1 {
		t.Errorf("%d coins in play, want %d", n, before+1)
	}
	d.Paint(8, 2, "stelo_poentoj") // later edits don't reach the running game
	if n := len(g.World.Room.All("stelo_poentoj")); n != before+1 {
		t.Errorf("the game shares the editor's level")
	}
}

// The strip's action buttons: undo is done by the screen, save and play are
// returned for the program.
func TestEditScreenActions(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := stelo.Files.ReadFile("levels/nivelo1.txt")
	l, _ := level.Parse(string(data))
	e := NewEditScreen(a, editor.New(l))
	orig := e.Doc.Text()
	e.Press(7*32+3, 2*32+3, false) // paint a block
	e.Release()
	if e.Doc.Text() == orig {
		t.Fatal("nothing painted")
	}
	if got := e.Press(actionX+3, paletteY+3, false); got != "undo" || e.Doc.Text() != orig {
		t.Errorf("undo button: %q, undone %v", got, e.Doc.Text() == orig)
	}
	if got := e.Press(actionX+editor.Cell+3, paletteY+3, false); got != "save" {
		t.Errorf("save button: %q", got)
	}
	if got := e.Press(actionX+2*editor.Cell+3, paletteY+3, false); got != "play" {
		t.Errorf("play button: %q", got)
	}
	if actionX < len(Palette)*editor.Cell {
		t.Errorf("the action buttons cover the palette")
	}
}
