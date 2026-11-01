package editor

import (
	"testing"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/level"
)

func load(t *testing.T, name string) *Doc {
	t.Helper()
	data, err := stelo.Files.ReadFile("levels/" + name + ".txt")
	if err != nil {
		t.Fatal(err)
	}
	l, err := level.Parse(string(data))
	if err != nil {
		t.Fatal(err)
	}
	return New(l)
}

// Painting, ladders, erasing, the unique player, and undo back to the
// original text; the edited level reads back as written.
func TestPaintUndoRoundTrip(t *testing.T) {
	d := load(t, "nivelo1")
	orig := d.Text()
	cx, cy := -1, -1 // an empty cell
	for y := 0; y*Cell < d.L.H && cx < 0; y++ {
		for x := 0; x*Cell < d.L.W; x++ {
			if o, l := d.At(x, y); o == "" && !l {
				cx, cy = x, y
				break
			}
		}
	}
	if !d.Paint(cx, cy, "block") {
		t.Fatal("painting a block changed nothing")
	}
	if o, _ := d.At(cx, cy); o != "block" {
		t.Fatalf("cell holds %q", o)
	}
	d.Paint(cx, cy, "monster") // replaces the block
	if o, _ := d.At(cx, cy); o != "monster" {
		t.Fatalf("replace: %q", o)
	}
	d.Paint(cx, cy, Ladder) // shares the cell
	if o, l := d.At(cx, cy); o != "monster" || !l {
		t.Fatalf("ladder: %q %v", o, l)
	}
	d.Paint(cx+1, cy, "character") // moves the player
	n := 0
	for _, in := range d.L.Instances {
		if in.Object == "character" {
			n++
		}
	}
	if n != 1 {
		t.Fatalf("%d players", n)
	}
	d.Paint(cx, cy, Erase)
	if o, l := d.At(cx, cy); o != "" || l {
		t.Fatalf("erase left %q %v", o, l)
	}
	edited := d.Text()
	back, err := level.Parse(edited)
	if err != nil || level.Format(back) != edited {
		t.Fatalf("the edited level doesn't read back (%v)", err)
	}
	for d.Undo() {
	}
	if d.Text() != orig {
		t.Error("undo didn't bring the original back")
	}
}

// Shrinking drops what falls outside, and undo restores it.
func TestResize(t *testing.T) {
	d := load(t, "nivelo2")
	orig := d.Text()
	if !d.Resize(20, 15) {
		t.Fatal("no resize")
	}
	for _, in := range d.L.Instances {
		if in.X >= 640 || in.Y >= 480 {
			t.Fatalf("%s at %d,%d kept outside", in.Object, in.X, in.Y)
		}
	}
	d.Undo()
	if d.Text() != orig {
		t.Error("undo didn't restore the size and the instances")
	}
}
