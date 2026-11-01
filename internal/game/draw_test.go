package game

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/LaPingvino/stelo"
)

// The first frame of level 1 matches gm53-reader's 1:1 render of the room
// (original/rooms/nivelo1.png), except where that render shows invisible
// objects (markers) and sprite-less ones (the life controller, drawn there as a
// placeholder), which the game does not draw.
func TestLevel1MatchesOriginalRender(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	w, err := LoadLevel(stelo.Files, a, "nivelo1")
	if err != nil {
		t.Fatal(err)
	}
	got := image.NewRGBA(image.Rect(0, 0, ScreenW, ScreenH))
	w.Room.G = nil // the level only, without the life controller's HUD
	w.Draw(got)
	f, err := os.Open("../../original/rooms/nivelo1.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	ref, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	skip := map[image.Point]bool{}
	for _, in := range w.Room.Instances {
		if !in.Obj.Visible || in.Obj.Sprite == "" {
			skip[image.Pt(int(in.X)/32, int(in.Y)/32)] = true
		}
	}
	bad := 0
	for y := 0; y < ScreenH; y++ {
		for x := 0; x < ScreenW; x++ {
			if skip[image.Pt(x/32, y/32)] {
				continue
			}
			r, g, b, _ := ref.At(x, y).RGBA()
			c := got.RGBAAt(x, y)
			if uint8(r>>8) != c.R || uint8(g>>8) != c.G || uint8(b>>8) != c.B {
				bad++
			}
		}
	}
	if bad != 0 {
		t.Errorf("%d pixels differ from the original render", bad)
	}
}
