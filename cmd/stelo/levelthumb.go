package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/level"
)

// The level pick's pictures (docs/UI.md §3: "a grid of the levels reached so
// far, each with a small picture of the level"): a level's map, small — its
// background colour, and each thing a mark in its colour.

var thumbColours = map[string]color.RGBA{
	"block": {0x9a, 0x86, 0x3c, 255}, "ladder": {0x3c, 0xa0, 0x4a, 255}, "stelo_poentoj": {0xff, 0xd8, 0x3a, 255},
	"monster": {0xe0, 0x4a, 0x2a, 255}, "b_monster": {0xe0, 0x4a, 0x2a, 255}, "flyer": {0x7a, 0x8a, 0x2a, 255},
	"death": {0xff, 0x8c, 0x1a, 255}, "character": {0xff, 0xff, 0xff, 255}, "levelexit": {0x3a, 0xd8, 0xe8, 255},
	"bonus": {0xff, 0x6a, 0x9a, 255}, "transirejo": {0xb0, 0x6a, 0xe0, 255}, "iru_al_bonuso": {0x9a, 0xe8, 0xff, 255},
}

// levelThumb draws level l into a w×h picture (the level fitted, centred).
func levelThumb(l *level.Level, w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{0x1a, 0x12, 0x0c, 255}), image.Point{}, draw.Src)
	if l == nil || l.W <= 0 || l.H <= 0 {
		return img
	}
	// The first screen: a level wider or taller than 640×480 shows the
	// screen around Stelo's start — what one sees on entering it.
	vx, vy, vw, vh := 0, 0, min(l.W, 640), min(l.H, 480)
	for _, in := range l.Instances {
		if in.Object == "character" {
			vx = max(0, min(in.X-vw/2, l.W-vw))
			vy = max(0, min(in.Y-vh/2, l.H-vh))
		}
	}
	k := min(float64(w)/float64(vw), float64(h)/float64(vh))
	ox, oy := (float64(w)-float64(vw)*k)/2, (float64(h)-float64(vh)*k)/2
	bg := color.RGBA{uint8(l.Colour >> 16), uint8(l.Colour >> 8), uint8(l.Colour), 255}
	if l.Background != "" { // the levels show their background picture (wood), not the colour
		bg = color.RGBA{0xc9, 0x9a, 0x62, 255}
	}
	draw.Draw(img, image.Rect(int(ox), int(oy), int(ox+float64(vw)*k), int(oy+float64(vh)*k)), image.NewUniform(bg), image.Point{}, draw.Src)
	cell := max(1, int(32*k+0.5))
	for _, in := range l.Instances {
		c, ok := thumbColours[in.Object]
		if !ok || in.X+32 <= vx || in.X >= vx+vw || in.Y+32 <= vy || in.Y >= vy+vh {
			continue
		}
		x, y := int(ox+float64(in.X-vx)*k), int(oy+float64(in.Y-vy)*k)
		draw.Draw(img, image.Rect(x, y, x+cell, y+cell), image.NewUniform(c), image.Point{}, draw.Src)
	}
	return img
}

// levelFor is level n of the chosen pack (nil when it can't be read).
func (a *app) levelFor(n int) *level.Level {
	if a.pack > 0 && a.pack < len(a.packs) {
		if ls := a.packs[a.pack].Levels; n >= 1 && n <= len(ls) {
			return ls[n-1]
		}
		return nil
	}
	return a.levelNamed(fmt.Sprintf("nivelo%d", n))
}

// levelNamed is the original's level of that room name (nil when it can't be read).
func (a *app) levelNamed(name string) *level.Level {
	data, err := stelo.Files.ReadFile("levels/" + name + ".txt")
	if err != nil {
		return nil
	}
	l, err := level.Parse(string(data))
	if err != nil {
		return nil
	}
	return l
}
