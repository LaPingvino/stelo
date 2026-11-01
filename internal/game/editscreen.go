package game

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/LaPingvino/stelo/internal/editor"
)

// EditScreen is the level editor's screen (LOOP item 10): the level drawn
// as the game draws it, the invisible objects (markers, spawners) shown
// too, and a palette strip of tools along the bottom.
type EditScreen struct {
	Assets *Assets
	Doc    *editor.Doc
	Tool   string
	ViewX  int
	ViewY  int
	// MouseX, MouseY is the pointer (screen pixels) for the cell cursor.
	MouseX, MouseY int

	world   *World
	dirty   bool
	painted map[[2]int]bool // cells painted in the current drag
}

// Palette is the tools in the strip's order: each object's own sprite
// shows it.
var Palette = []string{"block", editor.Ladder, "stelo_poentoj", "bonus", "death", "monster", "flyer", "marker",
	"transirejo", "iru_al_bonuso", "levelexit", "character", "b_monster", "monsterfarilo", "finmonsterfarilo", editor.Erase}

const paletteY = ScreenH - editor.Cell // the strip's top

// NewEditScreen edits a document with the block tool.
func NewEditScreen(a *Assets, d *editor.Doc) *EditScreen {
	return &EditScreen{Assets: a, Doc: d, Tool: "block", dirty: true}
}

// toolSprite is the sprite that shows a tool.
func (e *EditScreen) toolSprite(tool string) *Sprite {
	name := tool
	switch tool {
	case editor.Ladder:
		name = "ladder"
	case editor.Erase:
		return nil
	}
	if d := e.Assets.Objects[name]; d != nil {
		return e.Assets.Sprites[d.Sprite]
	}
	return nil
}

// Actions are the strip's buttons after the tools (for touch screens; the
// keys do the same: Ctrl+Z, Ctrl+S, P).
var Actions = []string{"open", "undo", "save", "play"}

// actionX is where the action buttons start: the strip's right end.
const actionX = ScreenW - 4*editor.Cell

// Press is a mouse press or a touch at screen (x, y): on the strip it picks
// a tool or returns an action ("undo" is done here; "save" and "play" are
// for the caller, as is "open"); on the level it starts painting (right =
// erase).
func (e *EditScreen) Press(x, y int, erase bool) (action string) {
	e.painted = map[[2]int]bool{}
	if y >= paletteY {
		e.painted = nil
		if x >= actionX {
			action = Actions[(x-actionX)/editor.Cell]
			if action == "undo" {
				e.Undo()
			}
			return action
		}
		if i := x / editor.Cell; i < len(Palette) {
			e.Tool = Palette[i]
		}
		return ""
	}
	e.Drag(x, y, erase)
	return ""
}

// Drag paints the cell under (x, y), once per cell per press.
func (e *EditScreen) Drag(x, y int, erase bool) {
	if e.painted == nil || y >= paletteY {
		return
	}
	cx, cy := (x+e.ViewX)/editor.Cell, (y+e.ViewY)/editor.Cell
	if e.painted[[2]int{cx, cy}] {
		return
	}
	e.painted[[2]int{cx, cy}] = true
	tool := e.Tool
	if erase {
		tool = editor.Erase
	}
	if e.Doc.Paint(cx, cy, tool) {
		e.dirty = true
	}
}

// Release ends the press.
func (e *EditScreen) Release() { e.painted = nil }

// Undo takes back the last change.
func (e *EditScreen) Undo() {
	if e.Doc.Undo() {
		e.dirty = true
	}
}

// Scroll moves the view by whole cells, inside the level (the strip covers
// the bottom row, so the view may go one row further).
func (e *EditScreen) Scroll(dx, dy int) {
	e.ViewX = max(0, min(e.Doc.L.W-ScreenW, e.ViewX+dx*editor.Cell))
	e.ViewY = max(0, min(e.Doc.L.H-paletteY, e.ViewY+dy*editor.Cell))
}

// Draw draws the level, the hidden objects, the cursor and the strip.
func (e *EditScreen) Draw(dst *image.RGBA) {
	if e.dirty || e.world == nil {
		e.world = NewWorld(e.Assets, e.Doc.L)
		e.dirty = false
	}
	w := e.world
	w.ViewX, w.ViewY = e.ViewX, e.ViewY
	w.Draw(dst)
	// Invisible objects (markers, spawners, the exit's helpers): drawn
	// half-transparent so they can be placed and seen.
	half := image.NewUniform(color.Alpha{128})
	for _, in := range w.Room.Instances {
		if in.Obj.Visible {
			continue
		}
		if sp := e.Assets.Sprites[in.Sprite]; sp != nil && len(sp.Frames) > 0 {
			img := sp.Frames[0]
			r := img.Bounds().Add(image.Pt(int(in.X)-e.ViewX, int(in.Y)-e.ViewY))
			draw.DrawMask(dst, r, img, image.Point{}, half, image.Point{}, draw.Over)
		}
	}
	// The cell under the pointer.
	if e.MouseY < paletteY {
		cx := (e.MouseX+e.ViewX)/editor.Cell*editor.Cell - e.ViewX
		cy := (e.MouseY+e.ViewY)/editor.Cell*editor.Cell - e.ViewY
		outline(dst, image.Rect(cx, cy, cx+editor.Cell, cy+editor.Cell), color.RGBA{255, 255, 255, 255})
	}
	// The strip: each tool's sprite, the chosen one outlined in yellow.
	draw.Draw(dst, image.Rect(0, paletteY, ScreenW, ScreenH), image.NewUniform(color.RGBA{40, 40, 48, 255}), image.Point{}, draw.Src)
	for i, tool := range Palette {
		x := i * editor.Cell
		if sp := e.toolSprite(tool); sp != nil && len(sp.Frames) > 0 {
			img := sp.Frames[0]
			draw.Draw(dst, img.Bounds().Add(image.Pt(x, paletteY)), img, image.Point{}, draw.Over)
		} else if tool == editor.Erase { // a red cross
			for d := 6; d < editor.Cell-6; d++ {
				dst.Set(x+d, paletteY+d, color.RGBA{220, 40, 40, 255})
				dst.Set(x+editor.Cell-1-d, paletteY+d, color.RGBA{220, 40, 40, 255})
			}
		}
		if tool == e.Tool {
			outline(dst, image.Rect(x, paletteY, x+editor.Cell, ScreenH), color.RGBA{255, 220, 0, 255})
		}
	}
	// The action buttons: undo (a hook arrow), save (a down arrow onto a
	// bar), play (a triangle) — drawn, the original has no art for them.
	white := color.RGBA{230, 230, 230, 255}
	for i, a := range Actions {
		x0, y0 := actionX+i*editor.Cell, paletteY
		outline(dst, image.Rect(x0+1, y0+1, x0+editor.Cell-1, ScreenH-1), color.RGBA{90, 90, 100, 255})
		switch a {
		case "open": // a folder: an outline with a tab
			for d := 0; d < 18; d++ {
				dst.Set(x0+7+d, y0+11, white)
				dst.Set(x0+7+d, y0+24, white)
			}
			for d := 0; d < 13; d++ {
				dst.Set(x0+7, y0+11+d, white)
				dst.Set(x0+24, y0+11+d, white)
			}
			for d := 0; d < 7; d++ {
				dst.Set(x0+7+d, y0+8, white)
			}
			dst.Set(x0+7, y0+9, white)
			dst.Set(x0+7, y0+10, white)
			dst.Set(x0+14, y0+9, white)
			dst.Set(x0+15, y0+10, white)
		case "undo":
			for d := 0; d < 12; d++ { // the arc's top and right side, then the arrowhead
				dst.Set(x0+10+d, y0+10, white)
				dst.Set(x0+22, y0+10+d, white)
				dst.Set(x0+10+d, y0+22, white)
			}
			for d := 0; d < 5; d++ {
				dst.Set(x0+10+d, y0+10-d, white)
				dst.Set(x0+10+d, y0+10+d, white)
			}
		case "save":
			for d := 0; d < 14; d++ {
				dst.Set(x0+16, y0+7+d, white)
				dst.Set(x0+9+d, y0+24, white)
			}
			for d := 0; d < 6; d++ {
				dst.Set(x0+16-d, y0+20-d, white)
				dst.Set(x0+16+d, y0+20-d, white)
			}
		case "play":
			for yy := 0; yy < 16; yy++ {
				w := 8 - abs(yy-8)
				for xx := 0; xx < w*2; xx++ {
					dst.Set(x0+11+xx, y0+8+yy, color.RGBA{80, 220, 80, 255})
				}
			}
		}
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func outline(dst *image.RGBA, r image.Rectangle, c color.Color) {
	for x := r.Min.X; x < r.Max.X; x++ {
		dst.Set(x, r.Min.Y, c)
		dst.Set(x, r.Max.Y-1, c)
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		dst.Set(r.Min.X, y, c)
		dst.Set(r.Max.X-1, y, c)
	}
}

// DrawPackChoice draws the title's pack line (when there are packs besides
// the original, as in The Puzzle Pits): "<  name  >" centred in the empty
// strip under the title's instructions (clear of the original's art), in the
// labels' font.
func (g *Game) DrawPackChoice(dst *image.RGBA, name string) {
	f := g.TextFont
	if f == nil {
		return
	}
	s := "<   " + name + "   >"
	w := font.MeasureString(f, s).Ceil()
	m := f.Metrics()
	y := 36
	d := font.Drawer{Dst: dst, Src: image.NewUniform(color.RGBA{255, 255, 255, 255}), Face: f, Dot: fixed.P((ScreenW-w)/2, y+m.Ascent.Ceil())}
	d.DrawString(s)
}

// Prompt is a question asked over the editor (a pack's name, its author):
// the label and the text typed so far.
type Prompt struct {
	Label, Text string
}

// DrawPrompt draws the prompt in a dark box across the middle of the screen.
func (e *EditScreen) DrawPrompt(dst *image.RGBA, p *Prompt, f font.Face) {
	if p == nil || f == nil {
		return
	}
	m := f.Metrics()
	h := 2*m.Height.Ceil() + 24
	y0 := (ScreenH - h) / 2
	r := image.Rect(40, y0, ScreenW-40, y0+h)
	draw.Draw(dst, r, image.NewUniform(color.RGBA{20, 20, 28, 240}), image.Point{}, draw.Over)
	outline(dst, r, color.RGBA{255, 220, 0, 255})
	white := image.NewUniform(color.RGBA{255, 255, 255, 255})
	d := font.Drawer{Dst: dst, Src: white, Face: f, Dot: fixed.P(r.Min.X+12, y0+10+m.Ascent.Ceil())}
	d.DrawString(p.Label)
	d = font.Drawer{Dst: dst, Src: image.NewUniform(color.RGBA{255, 220, 0, 255}), Face: f,
		Dot: fixed.P(r.Min.X+12, y0+14+m.Height.Ceil()+m.Ascent.Ceil())}
	d.DrawString(hats(p.Text) + "_")
}
