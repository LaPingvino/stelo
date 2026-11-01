package main

import (
	"image"
	"image/color"

	"github.com/LaPingvino/stelo/internal/gm"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Touch controls for phones (LOOP item 12): four arrow pads in the screen's
// bottom corners, shown once the screen has been touched. A finger on a pad
// holds that arrow key; a tap elsewhere counts as "any key" (keyboard_wait).

const padSize = 64

type pad struct {
	key  gm.Key
	r    image.Rectangle
	text string
}

// pads in game coordinates (the 640×480 screen): left/right bottom-left,
// up/down bottom-right.
var pads = []pad{
	{gm.KeyLeft, image.Rect(8, 480-8-padSize, 8+padSize, 480-8), "<"},
	{gm.KeyRight, image.Rect(16+padSize, 480-8-padSize, 16+2*padSize, 480-8), ">"},
	{gm.KeyUp, image.Rect(640-8-padSize, 480-16-2*padSize, 640-8, 480-16-padSize), "^"},
	{gm.KeyDown, image.Rect(640-8-padSize, 480-8-padSize, 640-8, 480-8), "v"},
}

type touchState struct {
	seen   bool
	held   map[gm.Key]bool
	tapped bool
}

// update reads the touches for this frame.
func (t *touchState) update() {
	t.held = map[gm.Key]bool{}
	t.tapped = false
	ids := ebiten.AppendTouchIDs(nil)
	if len(ids) > 0 {
		t.seen = true
	}
	for _, id := range ids {
		x, y := ebiten.TouchPosition(id)
		for _, p := range pads {
			if image.Pt(x, y).In(p.r) {
				t.held[p.key] = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		onPad := false
		for _, p := range pads {
			onPad = onPad || image.Pt(x, y).In(p.r)
		}
		if !onPad {
			t.tapped = true
		}
	}
}

// draw draws the pads, translucent, once touch has been used.
func (t *touchState) draw(screen *ebiten.Image) {
	if !t.seen {
		return
	}
	for _, p := range pads {
		a := uint8(70)
		if t.held[p.key] {
			a = 140
		}
		vector.FillRect(screen, float32(p.r.Min.X), float32(p.r.Min.Y), padSize, padSize, color.NRGBA{255, 255, 255, a}, false)
		drawArrow(screen, p)
	}
}

// drawArrow draws a filled triangle pointing the pad's way.
func drawArrow(screen *ebiten.Image, p pad) {
	cx, cy := float32(p.r.Min.X+padSize/2), float32(p.r.Min.Y+padSize/2)
	const s = 18
	var path vector.Path
	switch p.key {
	case gm.KeyLeft:
		path.MoveTo(cx-s, cy)
		path.LineTo(cx+s/2, cy-s)
		path.LineTo(cx+s/2, cy+s)
	case gm.KeyRight:
		path.MoveTo(cx+s, cy)
		path.LineTo(cx-s/2, cy-s)
		path.LineTo(cx-s/2, cy+s)
	case gm.KeyUp:
		path.MoveTo(cx, cy-s)
		path.LineTo(cx-s, cy+s/2)
		path.LineTo(cx+s, cy+s/2)
	case gm.KeyDown:
		path.MoveTo(cx, cy+s)
		path.LineTo(cx-s, cy-s/2)
		path.LineTo(cx+s, cy-s/2)
	}
	path.Close()
	op := &vector.FillOptions{}
	vector.FillPath(screen, &path, op, &vector.DrawPathOptions{ColorScale: func() ebiten.ColorScale {
		var c ebiten.ColorScale
		c.ScaleWithColor(color.NRGBA{0, 0, 0, 160})
		return c
	}()})
}
