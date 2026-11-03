package game

import (
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// The highscore screen (original/objects/plejbonaj_listo.gml, in finlisto):
//
//	Create:     the life controller goes, sound_loop(finavenko), font
//	            "Arial Rondigita MT Grasa" (Arial Rounded MT Bold) 20, green.
//	Room start: if score > highscore_value(10): a loop drawing the plejbonaj
//	            background and "Vi tiom bone ludis ke vi aperos en#la listo de
//	            la plej bonaj ludantoj!##Via nomo: "+keyboard_string+"_" at
//	            (70, 70) until Enter, then highscore_add(name, score); else the
//	            background and "Ho ve... Vi ne suficxe bone ludis por#aperi en
//	            la listo de la plej bonaj#ludantoj..." and keyboard_wait().
//	Draw:       draw_highscore(x+5, y+5, x+480, y+320); pluiru = 1.
//	Step:       if pluiru: keyboard_wait(); sound_stop_all(); game_restart().

// Highscore is one entry of the table.
type Highscore struct {
	Name  string
	Score int
}

// Highscores is Game Maker's ten-entry table, best first.
// ASSUMPTION [engine]: an empty table holds "<nobody>" with 0.
type Highscores [10]Highscore

// ClearedHighscores is highscore_clear().
func ClearedHighscores() Highscores {
	var h Highscores
	for i := range h {
		h[i] = Highscore{"<nobody>", 0}
	}
	return h
}

// Add is highscore_add: the entry goes in its place, the last one drops out.
func (h *Highscores) Add(name string, score int) {
	l := append(h[:], Highscore{name, score})
	sort.SliceStable(l, func(i, j int) bool { return l[i].Score > l[j].Score })
	copy(h[:], l)
}

// hsState is where the highscore screen is.
type hsState int

const (
	hsNone  hsState = iota
	hsName          // typing a name (until Enter)
	hsSorry         // "Ho ve…" until a key
	hsTable         // the table until a key, then the game restarts
)

const (
	hsGood   = "Vi tiom bone ludis ke vi aperos en#la listo de la plej bonaj ludantoj!##Via nomo: "
	hsSorry_ = "Ho ve... Vi ne suficxe bone ludis por#aperi en la listo de la plej bonaj#ludantoj..."
)

// startHighscores runs plejbonaj_listo's create and room start.
func (g *Game) startHighscores() {
	g.G.Vars["life_controller_gone"] = 1
	g.G.Play("@music:finavenko")
	if g.G.Score > g.Highscores[9].Score {
		g.G.KeyboardString = ""
		g.hs = hsName
	} else {
		g.hs = hsSorry
	}
}

// stepHighscores advances the screen; enter/anyKey are this step's keys.
// It returns true when the game should restart.
func (g *Game) stepHighscores(enter, anyKey bool) (restart bool) {
	switch g.hs {
	case hsName:
		if enter {
			g.Highscores.Add(g.G.KeyboardString, g.G.Score)
			g.hs = hsTable
			if g.SaveHighscores != nil {
				g.SaveHighscores(g.Highscores)
			}
		}
	case hsSorry:
		if anyKey {
			g.hs = hsTable
		}
	case hsTable:
		if anyKey {
			return true
		}
	}
	return false
}

// HighscoreFont is the screen's font: Liberation Sans Bold 20 pt standing in
// for Arial Rounded MT Bold.
func HighscoreFont(fsys fs.FS) (font.Face, error) {
	data, err := fs.ReadFile(fsys, "fonts/LiberationSans-Bold.ttf")
	if err != nil {
		return nil, err
	}
	f, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}
	return opentype.NewFace(f, &opentype.FaceOptions{Size: 20, DPI: 96, Hinting: font.HintingFull})
}

var hsGreen = image.NewUniform(color.RGBA{0, 255, 0, 255})

// drawGMText is draw_text: '#' breaks lines, the origin is the top-left.
func drawGMText(dst *image.RGBA, f font.Face, x, y int, s string) {
	m := f.Metrics()
	line := ""
	flush := func() {
		d := font.Drawer{Dst: dst, Src: hsGreen, Face: f, Dot: fixed.P(x, y+m.Ascent.Ceil())}
		d.DrawString(line)
		y += m.Height.Ceil()
		line = ""
	}
	for _, r := range hats(s) {
		if r == '#' {
			flush()
			continue
		}
		line += string(r)
	}
	flush()
}

// drawHighscores draws the highscore screen over the room.
func (g *Game) drawHighscores(dst *image.RGBA) {
	if g.hs != hsNone && g.World != nil && g.World.Overlay != nil {
		g.drawHighscoresUI()
		return
	}
	f := g.HSFont
	if f == nil {
		return
	}
	switch g.hs {
	case hsName, hsSorry:
		if bg := g.Assets.Backgrounds["plejbonaj"]; bg != nil {
			draw.Draw(dst, bg.Bounds(), bg, image.Point{}, draw.Src)
		}
		if g.hs == hsName {
			drawGMText(dst, f, 70, 70, Tr(hsGood)+g.G.KeyboardString+"_")
		} else {
			drawGMText(dst, f, 70, 70, Tr(hsSorry_))
		}
	case hsTable:
		// draw_highscore(x+5, y+5, x+480, y+320) for the instance at (64, 64).
		// ASSUMPTION [engine]: names at the left edge, scores right-aligned at
		// the right edge, the ten rows spread evenly over the height.
		x1, y1, x2, y2 := 69, 69, 544, 384
		m := f.Metrics()
		for i, e := range g.Highscores {
			y := y1 + i*(y2-y1)/10 + m.Ascent.Ceil()
			d := font.Drawer{Dst: dst, Src: hsGreen, Face: f, Dot: fixed.P(x1, y)}
			d.DrawString(e.Name)
			s := strconv.Itoa(e.Score)
			w := font.MeasureString(f, s).Ceil()
			d = font.Drawer{Dst: dst, Src: hsGreen, Face: f, Dot: fixed.P(x2-w, y)}
			d.DrawString(s)
		}
	}
}

// drawHighscoresUI is the new UI's highscore screen (docs/UI.md §5, step U5):
// the 2006 texts on a panel in Noto Sans, the name in a text field with a
// blinking caret, then the table with ranks.
func (g *Game) drawHighscoresUI() {
	o := g.World.Overlay
	panel, cream, gold := color.RGBA{0x2a, 0x1c, 0x12, 235}, color.RGBA{0xff, 0xf4, 0xd8, 255}, color.RGBA{0xff, 0xd0, 0x3a, 255}
	cx := float64(ScreenW) / 2
	switch g.hs {
	case hsName, hsSorry:
		o.AddRect(Rect{X: 0, Y: 0, W: ScreenW, H: ScreenH, Color: color.RGBA{0, 0, 0, 110}})
		o.AddRect(Rect{X: cx - 270, Y: 120, W: 540, H: 220, Color: panel})
		lines := strings.Split(hats(Tr(hsSorry_)), "#")
		if g.hs == hsName { // the 2006 text without its last line ("Via nomo: "), which the field takes
			lines = strings.Split(hats(Tr(hsGood)), "#")
			lines = lines[:len(lines)-2]
		}
		for i, line := range lines {
			o.Add(Text{S: line, X: cx, Y: 140 + float64(i)*28, Size: 18, Weight: Bold, Color: cream, Align: Centre})
		}
		if g.hs == hsName {
			o.Add(Text{S: hats(Tr("Via nomo:")), X: cx - 220, Y: 262, Size: 18, Weight: Bold, Color: gold})
			o.AddRect(Rect{X: cx - 110, Y: 254, W: 330, H: 40, Color: cream})
			caret := ""
			if g.World.Steps/15%2 == 0 {
				caret = "|"
			}
			o.Add(Text{S: g.G.KeyboardString + caret, X: cx - 100, Y: 262, Size: 18, Weight: Bold, Color: panel})
		}
	case hsTable:
		o.AddRect(Rect{X: cx - 210, Y: 40, W: 420, H: 400, Color: panel})
		o.Add(Text{S: Tr("Plej bonaj"), X: cx, Y: 54, Size: 30, Weight: Black, Color: gold, Align: Centre})
		for i, e := range g.Highscores {
			y := 112 + float64(i)*30
			col := cream
			if i == 0 {
				col = gold
			}
			o.Add(Text{S: strconv.Itoa(i+1) + ".", X: cx - 170, Y: y, Size: 18, Weight: Bold, Color: col, Align: Right})
			o.Add(Text{S: e.Name, X: cx - 155, Y: y, Size: 18, Weight: Bold, Color: col})
			o.Add(Text{S: strconv.Itoa(e.Score), X: cx + 180, Y: y, Size: 18, Weight: Bold, Color: col, Align: Right})
		}
	}
}

// ShowHighscores opens the highscore room on its table (the title's Plej
// bonaj, docs/UI.md §3): no name to enter, a key goes back.
func (g *Game) ShowHighscores() error {
	g.Title, g.waitRoom = false, false // as EnterLevel: past the 2006 title picture
	if err := g.enter("finlisto"); err != nil {
		return err
	}
	g.hs = hsTable
	return nil
}
