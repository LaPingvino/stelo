package main

import (
	"image"
	"image/color"

	"github.com/LaPingvino/stelo/internal/game"
	"github.com/LaPingvino/stelo/internal/gm"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Touch controls for phones (LOOP item 12, U0): arrow pads and a full-screen
// button, shown once the screen has been touched — in the room beside the game
// (a phone on its side: ← → in a column on the left, ↑ ↓ and full screen on
// the right) or under it (upright), never over the playfield (Joop 2026-11-02).
// They are laid out in screen pixels each frame; the frame is fitted into what
// is left. A finger on a pad holds that arrow key; a tap elsewhere counts as
// "any key" (keyboard_wait).

type pad struct {
	key gm.Key
	r   image.Rectangle // screen pixels
}

// touchLayout is where the frame and the controls go on a screen.
type touchLayout struct {
	area image.Rectangle // the frame is fitted into this
	pads []pad
	full image.Rectangle // the full-screen button
	snd  image.Rectangle // the sound on / off button
	size int             // a pad's side
}

// layoutTouch lays out the controls on a sw×sh screen (device pixels).
func layoutTouch(sw, sh int) touchLayout {
	d := ebiten.Monitor().DeviceScaleFactor()
	if d <= 0 {
		d = 1
	}
	m := int(10 * d)
	if sw > sh { // on its side: the controls in columns beside the frame
		p := min(int(64*d), (sh-3*m)/3)
		left := 2*p + 3*m
		right := p + 2*m
		y := sh - m - p
		return touchLayout{
			area: image.Rect(left, 0, sw-right, sh),
			pads: []pad{
				{gm.KeyLeft, image.Rect(m, y, m+p, y+p)},
				{gm.KeyRight, image.Rect(2*m+p, y, 2*m+2*p, y+p)},
				{gm.KeyUp, image.Rect(sw-m-p, y-m-p, sw-m, y-m)},
				{gm.KeyDown, image.Rect(sw-m-p, y, sw-m, y+p)},
			},
			full: image.Rect(m, m, m+p, m+p*2/3), // top left (the page's Klavaro button is top right)
			snd:  image.Rect(m, 2*m+p*2/3, m+p, 2*m+p*4/3),
			size: p,
		}
	}
	p := min(int(64*d), (sw-5*m)/4)
	fh := sw * game.ScreenH / game.ScreenW // the frame as wide as the screen
	y := min(sh-m-p, fh+m)                 // the controls right under it
	return touchLayout{                    // upright: the controls in a row under the frame
		area: image.Rect(0, 0, sw, y-m),
		pads: []pad{
			{gm.KeyLeft, image.Rect(m, y, m+p, y+p)},
			{gm.KeyRight, image.Rect(2*m+p, y, 2*m+2*p, y+p)},
			{gm.KeyUp, image.Rect(sw-2*m-2*p, y, sw-2*m-p, y+p)},
			{gm.KeyDown, image.Rect(sw-m-p, y, sw-m, y+p)},
		},
		full: image.Rect(sw/2-p-m/2, y+p/6, sw/2-m/2, y+p*5/6),
		snd:  image.Rect(sw/2+m/2, y+p/6, sw/2+p+m/2, y+p*5/6),
		size: p,
	}
}

type touchState struct {
	seen   bool
	held   map[gm.Key]bool
	tapped bool
	lay    touchLayout // last frame's layout (the touches are read against it)
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
		pt := image.Pt(ebiten.TouchPosition(id))
		for _, p := range t.lay.pads {
			if pt.In(p.r) {
				t.held[p.key] = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		pt := image.Pt(ebiten.TouchPosition(id))
		if pt.In(t.lay.full) {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
			continue
		}
		if pt.In(t.lay.snd) {
			toggleSound()
			continue
		}
		onPad := false
		for _, p := range t.lay.pads {
			onPad = onPad || pt.In(p.r)
		}
		if !onPad {
			t.tapped = true
		}
	}
}

// draw draws the pads and the full-screen button on the screen (outside the
// frame), once touch has been used.
func (t *touchState) draw(screen *ebiten.Image) {
	drawSoundButton(screen, t.lay.snd)
	if !t.seen {
		drawFullButton(screen, t.lay.full)
		return
	}
	for _, p := range t.lay.pads {
		a := uint8(70)
		if t.held[p.key] {
			a = 150
		}
		vector.FillRect(screen, float32(p.r.Min.X), float32(p.r.Min.Y), float32(p.r.Dx()), float32(p.r.Dy()), color.NRGBA{255, 255, 255, a}, false)
		drawArrow(screen, p)
	}
	drawFullButton(screen, t.lay.full)
}

// drawFullButton draws the full-screen button: four corner marks.
func drawFullButton(screen *ebiten.Image, f image.Rectangle) {
	if f.Empty() {
		return
	}
	vector.FillRect(screen, float32(f.Min.X), float32(f.Min.Y), float32(f.Dx()), float32(f.Dy()), color.NRGBA{255, 255, 255, 50}, false)
	// four corner marks: the full-screen sign
	c, w := float32(f.Dy())/4, float32(max(2, f.Dy()/16))
	x0, y0, x1, y1 := float32(f.Min.X)+c, float32(f.Min.Y)+c/1.2, float32(f.Max.X)-c, float32(f.Max.Y)-c/1.2
	col := color.NRGBA{255, 255, 255, 200}
	for _, l := range [][4]float32{{x0, y0, x0 + c, y0}, {x0, y0, x0, y0 + c}, {x1, y0, x1 - c, y0}, {x1, y0, x1, y0 + c},
		{x0, y1, x0 + c, y1}, {x0, y1, x0, y1 - c}, {x1, y1, x1 - c, y1}, {x1, y1, x1, y1 - c}} {
		vector.StrokeLine(screen, l[0], l[1], l[2], l[3], w, col, false)
	}
}

// drawArrow draws a filled triangle pointing the pad's way.
func drawArrow(screen *ebiten.Image, p pad) {
	cx, cy := float32(p.r.Min.X+p.r.Dx()/2), float32(p.r.Min.Y+p.r.Dy()/2)
	s := float32(p.r.Dx()) * 0.28
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

// phoneKeys is what the page's phone keyboard (web/index.html) sent since
// the last frame: characters ('\b' for Backspace) and Enter.
type phoneKeys struct {
	chars []rune
	enter bool
}

// drawSoundButton draws the sound button: a loudspeaker, with waves when the
// sound is on and a cross when it is off.
func drawSoundButton(screen *ebiten.Image, r image.Rectangle) {
	if r.Empty() {
		return
	}
	vector.FillRect(screen, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), color.NRGBA{255, 255, 255, 50}, false)
	h := float32(r.Dy())
	cx, cy := float32(r.Min.X)+float32(r.Dx())/2, float32(r.Min.Y)+h/2
	u := h / 8
	col := color.NRGBA{255, 255, 255, 210}
	var p vector.Path // the loudspeaker
	p.MoveTo(cx-3*u, cy-u)
	p.LineTo(cx-1.6*u, cy-u)
	p.LineTo(cx, cy-2.4*u)
	p.LineTo(cx, cy+2.4*u)
	p.LineTo(cx-1.6*u, cy+u)
	p.LineTo(cx-3*u, cy+u)
	p.Close()
	vector.FillPath(screen, &p, &vector.FillOptions{}, &vector.DrawPathOptions{ColorScale: func() ebiten.ColorScale {
		var c ebiten.ColorScale
		c.ScaleWithColor(col)
		return c
	}()})
	w := max(1.5, u/2.5)
	if soundMuted() {
		vector.StrokeLine(screen, cx+u, cy-1.4*u, cx+3.4*u, cy+1.4*u, w, col, true)
		vector.StrokeLine(screen, cx+u, cy+1.4*u, cx+3.4*u, cy-1.4*u, w, col, true)
		return
	}
	vector.StrokeLine(screen, cx+1.2*u, cy-1.2*u, cx+1.2*u, cy+1.2*u, w, col, true)
	vector.StrokeLine(screen, cx+2.4*u, cy-2*u, cx+2.4*u, cy+2*u, w, col, true)
}
