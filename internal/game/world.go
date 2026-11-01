// Package game is the Stelo remake's game: the world of one room and its
// drawing. Rules: docs/STELO-RULES.md.
package game

import (
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"sort"

	"github.com/LaPingvino/stelo/internal/level"
)

// Screen size (the original's window and view).
const ScreenW, ScreenH = 640, 480

// Instance is one object in the room.
type Instance struct {
	Def    *ObjectDef
	X, Y   float64
	Sprite string
	Frame  int
}

// World is the running room.
type World struct {
	Assets       *Assets
	Level        *level.Level
	Instances    []*Instance
	ViewX, ViewY int
	follow       *Instance
}

// LoadLevel reads levels/<name>.txt and places its instances.
func LoadLevel(fsys fs.FS, a *Assets, name string) (*World, error) {
	data, err := fs.ReadFile(fsys, "levels/"+name+".txt")
	if err != nil {
		return nil, err
	}
	l, err := level.Parse(string(data))
	if err != nil {
		return nil, err
	}
	w := &World{Assets: a, Level: l}
	for _, in := range l.Instances {
		def := a.Objects[in.Object]
		if def == nil {
			def = &ObjectDef{Name: in.Object}
		}
		inst := &Instance{Def: def, X: float64(in.X), Y: float64(in.Y), Sprite: def.Sprite}
		w.Instances = append(w.Instances, inst)
		if l.View != nil && l.View.Follow == in.Object && w.follow == nil {
			w.follow = inst
		}
	}
	if l.View != nil {
		w.ViewX, w.ViewY = l.View.X, l.View.Y
	}
	return w, nil
}

// UpdateView moves the view after the followed instance, as Game Maker does:
// when it is nearer than the border to an edge, the view moves to bring it
// back to the border, by at most the view speed per step (−1: no limit), and
// stays inside the room.
func (w *World) UpdateView() {
	v := w.Level.View
	if v == nil || w.follow == nil {
		return
	}
	step := func(pos, size, room, objPos, border, speed int) int {
		target := pos
		if objPos-pos < border {
			target = objPos - border
		} else if pos+size-objPos < border {
			target = objPos + border - size
		}
		d := target - pos
		if speed >= 0 {
			d = max(-speed, min(speed, d))
		}
		return max(0, min(room-size, pos+d))
	}
	w.ViewX = step(w.ViewX, v.W, w.Level.W, int(w.follow.X), v.BorderH, v.SpeedH)
	w.ViewY = step(w.ViewY, v.H, w.Level.H, int(w.follow.Y), v.BorderV, v.SpeedV)
}

// Draw renders the view into dst (ScreenW×ScreenH): the background colour,
// the tiled background, then the visible instances from the highest depth to
// the lowest.
func (w *World) Draw(dst *image.RGBA) {
	c := w.Level.Colour
	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.RGBA{uint8(c >> 16), uint8(c >> 8), uint8(c), 255}}, image.Point{}, draw.Src)
	if bg := w.Assets.Backgrounds[w.Level.Background]; bg != nil {
		bw, bh := bg.Bounds().Dx(), bg.Bounds().Dy()
		for y := -(w.ViewY % bh); y < ScreenH; y += bh {
			for x := -(w.ViewX % bw); x < ScreenW; x += bw {
				draw.Draw(dst, image.Rect(x, y, x+bw, y+bh), bg, image.Point{}, draw.Over)
			}
		}
	}
	insts := append([]*Instance(nil), w.Instances...)
	sort.SliceStable(insts, func(i, j int) bool { return insts[i].Def.Depth > insts[j].Def.Depth })
	for _, in := range insts {
		sp := w.Assets.Sprites[in.Sprite]
		if !in.Def.Visible || sp == nil || len(sp.Frames) == 0 {
			continue
		}
		img := sp.Frames[in.Frame%len(sp.Frames)]
		x, y := int(in.X)-w.ViewX, int(in.Y)-w.ViewY
		r := img.Bounds().Add(image.Pt(x, y))
		if r.Overlaps(dst.Bounds()) {
			draw.Draw(dst, r, img, image.Point{}, draw.Over)
		}
	}
}
