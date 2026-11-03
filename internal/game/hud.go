package game

import (
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"strconv"
	"strings"

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
// The life controller is persistent and placed once, in level 1; danklisto's
// creation code and plejbonaj_listo (finlisto) destroy it. The HUD is drawn in
// rooms with a player while it exists (not on the rules screen, klarigo,
// whose player is the legend's picture).

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

// TextFont loads tekstilo's font: action_draw_font("Arial Unicode MS", 11,
// green, bold) — Liberation Sans Bold standing in.
func TextFont(fsys fs.FS) (font.Face, error) {
	data, err := fs.ReadFile(fsys, "fonts/LiberationSans-Bold.ttf")
	if err != nil {
		return nil, err
	}
	f, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(f, &opentype.FaceOptions{Size: 11, DPI: 96, Hinting: font.HintingFull})
}

// drawTexts is tekstilo's Draw event: its text (creation code teksto='…') at
// (x+3, y) relative to the view, green (65280 = 0x00FF00), '#' a line break.
func (w *World) drawTexts(dst *image.RGBA) {
	if w.TextFont == nil {
		return
	}
	green := image.NewUniform(color.RGBA{0, 255, 0, 255})
	m := w.TextFont.Metrics()
	for _, in := range w.Room.All("tekstilo") {
		t, _ := in.VarsStr["teksto"]
		y := int(in.Y) - w.ViewY + m.Ascent.Ceil()
		for _, line := range strings.Split(hats(Tr(t)), "#") {
			d := font.Drawer{Dst: dst, Src: green, Face: w.TextFont, Dot: fixed.P(int(in.X)+3-w.ViewX, y)}
			d.DrawString(line)
			y += m.Height.Ceil()
		}
	}
}

// drawHUD draws the life controller's HUD for the room's player.
func (w *World) drawHUD(dst *image.RGBA, score, lives int) {
	players := w.Room.All("character")
	if w.Room.Paused != nil { // during a pause, where the player was when it began
		players = players[:0:0]
		for _, in := range w.Room.Paused {
			if in.Obj.IsA("character") && !in.Destroyed {
				players = append(players, in)
			}
		}
	}
	v := w.Room.G.Vars
	if len(players) == 0 || v["life_controller_alive"] != 1 || v["life_controller_gone"] == 1 {
		return
	}
	if w.Overlay != nil { // the new UI's slim bar (docs/UI.md §5)
		w.drawBar(dst, score, lives)
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
		if w.Skin != nil && len(w.Skin.Sprites["sprite_life"]) > 0 { // the skin's icon
			img = w.Skin.Sprites["sprite_life"][0]
		}
		sw := img.Bounds().Dx()
		for i := range max(lives, 0) {
			x, y := int(p.X)-w.ViewX+i*sw, int(p.Y)-10-w.ViewY
			draw.Draw(dst, img.Bounds().Add(image.Pt(x, y)), img, image.Point{}, draw.Over)
		}
	}
	// draw_score: caption and score, top-left aligned at (10, 5).
	if w.Overlay != nil {
		w.Overlay.Add(Text{S: Tr("Poentoj: ") + strconv.Itoa(score), X: 10, Y: 6, Size: 13, Weight: Bold, Color: color.Black})
	} else if w.Font != nil {
		d := font.Drawer{Dst: dst, Src: image.NewUniform(color.Black), Face: w.Font}
		d.Dot = fixed.P(10, 5+w.Font.Metrics().Ascent.Ceil())
		d.DrawString(Tr("Poentoj: ") + strconv.Itoa(score))
	}
}

// drawBar is the new UI's HUD (docs/UI.md §5, a remake change): one slim bar
// along the top — the score, the level's name, the lives as icons — instead
// of the 2006 cream box and the lives over the player's head; and the level's
// name large for two seconds when it starts.
func (w *World) drawBar(dst *image.RGBA, score, lives int) {
	const h = 22
	bar := image.Rect(0, 0, ScreenW, h)
	draw.Draw(dst, bar, image.NewUniform(color.RGBA{0x1c, 0x12, 0x0a, 170}), image.Point{}, draw.Over)
	cream := color.RGBA{0xff, 0xf4, 0xd8, 255}
	w.Overlay.Add(Text{S: Tr("Poentoj") + "  " + strconv.Itoa(score), X: 10, Y: 3, Size: 13, Weight: Bold, Color: cream})
	if w.Label != "" {
		w.Overlay.Add(Text{S: w.Label, X: ScreenW / 2, Y: 3, Size: 13, Weight: Bold, Color: cream, Align: Centre})
	}
	if sp := w.Assets.Sprites["sprite_life"]; sp != nil && len(sp.Frames) > 0 {
		img := sp.Frames[0]
		if w.Skin != nil && len(w.Skin.Sprites["sprite_life"]) > 0 {
			img = w.Skin.Sprites["sprite_life"][0]
		}
		iw, ih := img.Bounds().Dx(), img.Bounds().Dy()
		n := max(lives, 0)
		shown := min(n, 6)
		x := ScreenW - 10 - shown*(iw+3)
		if n > shown { // many lives: the icons and the number
			w.Overlay.Add(Text{S: "×" + strconv.Itoa(n), X: float64(x - 4), Y: 3, Size: 13, Weight: Bold, Color: cream, Align: Right})
			shown = 1
			x = ScreenW - 10 - (iw + 3)
		}
		for i := range shown {
			draw.Draw(dst, img.Bounds().Add(image.Pt(x+i*(iw+3), (h-ih)/2)), img, image.Point{}, draw.Over)
		}
	}
	if st := w.Steps; w.Label != "" && st < 60 && !w.HideBanner { // two seconds at 30 steps a second
		a := uint8(255)
		if st > 45 { // fading out over the last half second
			a = uint8(255 * (60 - st) / 15)
		}
		w.Overlay.Add(Text{S: w.Label, X: ScreenW/2 + 3, Y: ScreenH/2 - 37, Size: 48, Weight: Black, Color: color.NRGBA{0x3b, 0x22, 0x0c, a}, Align: Centre})
		w.Overlay.Add(Text{S: w.Label, X: ScreenW / 2, Y: ScreenH/2 - 40, Size: 48, Weight: Black, Color: color.NRGBA{0xff, 0xd0, 0x3a, a}, Align: Centre})
	}
}
