package game

import (
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"sort"
	"strconv"

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
			drawGMText(dst, f, 70, 70, hsGood+g.G.KeyboardString+"_")
		} else {
			drawGMText(dst, f, 70, 70, hsSorry_)
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
