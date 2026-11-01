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
