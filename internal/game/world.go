// Package game is the Stelo remake's game: a level run on the Game Maker core
// (internal/gm) and drawn as the original. Rules: docs/STELO-RULES.md.
package game

import (
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"sort"
	"strings"

	"github.com/LaPingvino/stelo/internal/gm"
	"github.com/LaPingvino/stelo/internal/level"
	"golang.org/x/image/font"
)

// Screen size (the original's window and view).
const ScreenW, ScreenH = 640, 480

// World is the running room.
type World struct {
	Assets       *Assets
	Level        *level.Level
	Room         *gm.Room
	ViewX, ViewY int
	follow       *gm.Instance
	// Skin, if set, replaces the original sprites when drawing (item 11).
	Skin *Skin
	// Font is the HUD font (HUDFont); nil = no score text. TextFont is the
	// text labels' font (TextFont).
	Font, TextFont font.Face
	// Tick counts the steps (the skin's animation clock).
	Tick int
	// Overlay, when set, receives the HUD's text to be drawn sharp at the
	// screen's resolution (docs/UI.md); nil: drawn into the frame as before.
	Overlay *Overlay
	// Label is the level's name for the HUD ("Nivelo 2", "Bonusnivelo"),
	// shown large for its first 60 steps; Steps counts the room's steps.
	Label string
	Steps int
	// HideBanner hides the level's name (the pause menu is open: the steps
	// stand still, so it would stay).
	HideBanner bool
}

// Object returns the engine object for an original object (sprite box,
// solidity, drawing defaults), shared by all its instances.
func (a *Assets) Object(name string) *gm.Object {
	if o := a.gmObjects[name]; o != nil {
		return o
	}
	o := &gm.Object{Name: name}
	if d := a.Objects[name]; d != nil {
		o.Solid, o.Sprite, o.Visible, o.Depth = d.Solid, d.Sprite, d.Visible, d.Depth
		if sp := a.Sprites[d.Sprite]; sp != nil {
			o.Mask = true
			o.BBox = gm.Box{L: sp.BBox[0], R: sp.BBox[1], T: sp.BBox[2], B: sp.BBox[3]}
		}
	}
	if a.gmObjects == nil {
		a.gmObjects = map[string]*gm.Object{}
	}
	a.gmObjects[name] = o
	if ev := objectEvents[name]; ev != nil {
		ev(a, o)
	}
	return o
}

// LoadLevel reads levels/<name>.txt and places its instances, running their
// create events in the level's order.
func LoadLevel(fsys fs.FS, a *Assets, name string) (*World, error) {
	data, err := fs.ReadFile(fsys, "levels/"+name+".txt")
	if err != nil {
		return nil, err
	}
	l, err := level.Parse(string(data))
	if err != nil {
		return nil, err
	}
	return NewWorld(a, l), nil
}

// NewWorld makes the running room for a level (its instances created in
// order, their creation code run).
func NewWorld(a *Assets, l *level.Level) *World {
	w := &World{Assets: a, Level: l, Room: &gm.Room{Name: l.Name, W: l.W, H: l.H, G: &gm.Globals{Lives: 3, Vars: map[string]float64{}}}}
	for _, in := range l.Instances {
		inst := w.Room.Create(a.Object(in.Object), float64(in.X), float64(in.Y))
		// Instance creation code: the only kind the originals use sets a
		// label's text (teksto='…';).
		if t, ok := strings.CutPrefix(in.Code, "teksto='"); ok {
			inst.VarsStr["teksto"] = strings.TrimSuffix(t, "';")
		}
		if l.View != nil && l.View.Follow == in.Object && w.follow == nil {
			w.follow = inst
		}
	}
	if l.View != nil {
		w.ViewX, w.ViewY = l.View.X, l.View.Y
	}
	return w
}

// Step runs one step of the room, then moves the view.
func (w *World) Step(in gm.Input) {
	w.Tick++ // the skin's animation clock (it runs through pauses too)
	w.Room.Step(in)
	w.UpdateView()
}

// UpdateView moves the view after the followed instance, as Game Maker does:
// when it is nearer than the border to an edge, the view moves to bring it
// back to the border, by at most the view speed per step (−1: no limit), and
// stays inside the room.
func (w *World) UpdateView() {
	v := w.Level.View
	if v == nil || w.follow == nil || w.Room.Paused != nil { // a pause shows the moment it began: the view waits too
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
	insts := append([]*gm.Instance(nil), w.Room.Drawn()...) // during a pause, the moment it began
	sort.SliceStable(insts, func(i, j int) bool { return insts[i].Obj.Depth > insts[j].Obj.Depth })
	var blocks map[[2]int]bool // the skin's floor tops: blocks with no block above
	if w.Skin != nil && w.Skin.BlockTop != nil {
		blocks = map[[2]int]bool{}
		for _, in := range w.Room.Drawn() {
			if in.Obj.Name == "block" && !in.Destroyed {
				blocks[[2]int{int(in.X), int(in.Y)}] = true
			}
		}
	}
	for _, in := range insts {
		sp := w.Assets.Sprites[in.Sprite]
		if !in.Obj.Visible || in.Destroyed || sp == nil || len(sp.Frames) == 0 {
			continue
		}
		frames := sp.Frames
		var img *image.RGBA
		if w.Skin != nil {
			skinned := false
			if fr := w.Skin.Sprites[in.Sprite]; len(fr) > 0 {
				frames, skinned = fr, true
			}
			if in.Obj.Name == "block" && blocks != nil && !blocks[[2]int{int(in.X), int(in.Y) - 32}] {
				frames, skinned = w.Skin.BlockTop, true
			}
			if skinned { // the skin's own, calmer pace
				img = w.Skin.Frame(in.Sprite, frames, w.Tick, in.XStart, in.YStart, in.X, in.Y, in.XPrevious, in.YPrevious)
			}
		}
		if img == nil {
			img = frames[int(in.ImageIndex)%len(frames)]
		}
		x, y := int(in.X)-w.ViewX, int(in.Y)-w.ViewY
		r := img.Bounds().Add(image.Pt(x, y))
		if r.Overlaps(dst.Bounds()) {
			draw.Draw(dst, r, img, image.Point{}, draw.Over)
		}
	}
	w.drawTexts(dst)
	if g := w.Room.ShownG(); g != nil {
		w.drawHUD(dst, g.Score, g.Lives)
	}
}

// DrawTitle draws the title picture (background stelo_komenco) full-screen.
func (w *World) DrawTitle(dst *image.RGBA) {
	if bg := w.Assets.Backgrounds["stelo_komenco"]; bg != nil {
		draw.Draw(dst, dst.Bounds(), bg, image.Point{}, draw.Src)
	}
}

// Clone copies the world's running state (gm.Room.Clone), the view and the
// instance it follows; assets and fonts are shared.
func (w *World) Clone() *World {
	c := *w
	var m map[*gm.Instance]*gm.Instance
	c.Room, m = w.Room.Clone()
	c.follow = m[w.follow]
	return &c
}
