package game

import (
	"image"
	"image/color"
	"image/draw"

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

// Press is a mouse press or a touch at screen (x, y): on the strip it picks
// a tool; on the level it starts painting (right = erase).
func (e *EditScreen) Press(x, y int, erase bool) {
	e.painted = map[[2]int]bool{}
	if y >= paletteY {
		if i := x / editor.Cell; i < len(Palette) {
			e.Tool = Palette[i]
		}
		e.painted = nil
		return
	}
	e.Drag(x, y, erase)
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
