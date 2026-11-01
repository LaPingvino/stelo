package game

import (
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// The life controller's drawing (original/objects/life_controller.gml, Draw):
//
//	action_color(11855591, 0)          brush (fill) light yellow, pen (outline) black
//	action_draw_rectangle(view_left+5, view_top+5, view_left+100, view_top+25)
//	action_draw_life_images(character.x, character.y-10, sprite_life)
//	action_draw_score(view_left+10, view_top+5, 'Poentoj: ')
//
// in the font set at its first room start: action_draw_font("@Arial Unicode
// MS", 10, black). Liberation Sans (metric-compatible with Arial) stands in;
// Game Maker sizes are points at 96 dpi.
//
// The life controller is persistent and placed once, in level 1; it is gone in
// danklisto and finlisto. The HUD is drawn in every room with a player, which
// is the same set of rooms.

var hudFill = color.RGBA{231, 230, 180, 255} // 11855591 = 0xB4E6E7 (BGR)

// HUDFont loads the stand-in font at the original's size.
func HUDFont(fsys fs.FS) (font.Face, error) {
	data, err := fs.ReadFile(fsys, "fonts/LiberationSans-Regular.ttf")
	if err != nil {
		return nil, err
	}
	f, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(f, &opentype.FaceOptions{Size: 10, DPI: 96, Hinting: font.HintingFull})
}

// drawHUD draws the life controller's HUD for the room's player.
func (w *World) drawHUD(dst *image.RGBA, score, lives int) {
	players := w.Room.All("character")
	if len(players) == 0 {
		return
	}
	p := players[0]
	// draw_rectangle: brush fill, then the pen outline, corners inclusive.
	x1, y1, x2, y2 := 5, 5, 100, 25
	draw.Draw(dst, image.Rect(x1, y1, x2+1, y2+1), image.NewUniform(color.Black), image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(x1+1, y1+1, x2, y2), image.NewUniform(hudFill), image.Point{}, draw.Src)
	// draw_life_images: one sprite_life per life, side by side.
	if sp := w.Assets.Sprites["sprite_life"]; sp != nil && len(sp.Frames) > 0 {
		img := sp.Frames[0]
		sw := img.Bounds().Dx()
		for i := range max(lives, 0) {
			x, y := int(p.X)-w.ViewX+i*sw, int(p.Y)-10-w.ViewY
			draw.Draw(dst, img.Bounds().Add(image.Pt(x, y)), img, image.Point{}, draw.Over)
		}
	}
	// draw_score: caption and score, top-left aligned at (10, 5).
	if w.Font != nil {
		d := font.Drawer{Dst: dst, Src: image.NewUniform(color.Black), Face: w.Font}
		d.Dot = fixed.P(10, 5+w.Font.Metrics().Ascent.Ceil())
		d.DrawString("Poentoj: " + strconv.Itoa(score))
	}
}
