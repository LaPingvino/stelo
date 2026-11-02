package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/LaPingvino/stelo/internal/game"
)

// The game's 640×480 frame on the real screen (docs/UI.md §1): scaled by
// whole steps and never smoothed, so the pixel art stays sharp — a black
// border where the screen doesn't divide evenly — as long as whole steps fill
// at least 90% of the room; otherwise (a dense phone screen, an odd window)
// by the largest fraction that fits, still unsmoothed. Layout gives
// the screen's real resolution, so what is drawn at it (the menus, later) is
// sharp too.

// frameView is where the frame lies on the screen: its scale and top-left.
type frameView struct {
	scale  float64
	ox, oy float64
}

var view = frameView{scale: 1}

// fitIn places the frame in an area of the screen.
func fitIn(r image.Rectangle) frameView {
	v := fit(r.Dx(), r.Dy())
	v.ox += float64(r.Min.X)
	v.oy += float64(r.Min.Y)
	return v
}

// fit places the frame on a sw×sh screen.
func fit(sw, sh int) frameView {
	k := min(float64(sw)/game.ScreenW, float64(sh)/game.ScreenH)
	if w := float64(int(k)); w >= 1 && w/k >= 0.9 {
		k = w // whole steps, every pixel equally wide, when they waste little room
	}
	return frameView{scale: k, ox: (float64(sw) - game.ScreenW*k) / 2, oy: (float64(sh) - game.ScreenH*k) / 2}
}

// toFrame turns a screen position (mouse, touch) into frame coordinates.
func toFrame(x, y int) (int, int) {
	return int((float64(x) - view.ox) / view.scale), int((float64(y) - view.oy) / view.scale)
}

// present draws the frame (and what was drawn onto it) onto the screen, in
// area (the whole screen, or what the touch controls leave).
func present(screen, frame *ebiten.Image, area image.Rectangle) {
	view = fitIn(area)
	screen.Fill(color.Black)
	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
	op.GeoM.Scale(view.scale, view.scale)
	op.GeoM.Translate(view.ox, view.oy)
	screen.DrawImage(frame, op)
}

// screenSize is the real resolution for a window of outW×outH (device pixels).
func screenSize(outW, outH int) (int, int) {
	s := ebiten.Monitor().DeviceScaleFactor()
	return int(float64(outW) * s), int(float64(outH) * s)
}
