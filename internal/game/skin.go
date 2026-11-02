package game

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/fs"
)

// Skin is a drawing-only replacement of the original sprites (LOOP item 11):
// frames by original sprite name. Collision boxes stay the original's, so
// the rules don't change. Sprites it doesn't name keep the original art.
type Skin struct {
	Sprites map[string][]*image.RGBA
	// Rate is each sprite's animation speed in frames per step (the skin's
	// own pace, slower than the original's image_speed: at 30 steps a second
	// 0.2 is 6 frames a second); sprites it doesn't name use defaultRate.
	Rate map[string]float64
	// IdleFirst sprites show their first frame while the instance stands
	// still (the protagonist only walks its cycle when it moves).
	IdleFirst map[string]bool
	// BlockTop is drawn for a block with no block right above it (a floor
	// with grass), the plain block frames elsewhere.
	BlockTop []*image.RGBA
}

// LoadNovaSkin builds the skin from Joop's commissioned 32×32 set
// (art/nova/). The roles, ours to choose (docs/DESIGN.md "New art"):
// the protagonist star is the player (its left sheet facing left; climbing
// shows its first frame; the life icon it, shrunk), the orange villain the walking monster and the
// bonus level's (mirrored to face left), the green villain the flyer (it
// can't be stomped, so it looks meaner), the checkpoint's star a star coin, the
// heart and diamond themselves, the checkpoint unused (the orb on a pole) and
// used (the shining star), the spikes the fire, the arrow sign the level exit, the
// ladder, and wall / floor-with-grass for blocks.
func LoadNovaSkin(fsys fs.FS) (*Skin, error) {
	sheet := func(name string) ([]*image.RGBA, error) {
		f, err := fsys.Open("art/nova/" + name + ".png")
		if err != nil {
			return nil, err
		}
		defer f.Close()
		img, err := png.Decode(f)
		if err != nil {
			return nil, err
		}
		var frames []*image.RGBA
		b := img.Bounds()
		for x := b.Min.X; x+32 <= b.Max.X; x += 32 {
			fr := image.NewRGBA(image.Rect(0, 0, 32, 32))
			draw.Draw(fr, fr.Bounds(), img, image.Pt(x, b.Min.Y), draw.Src)
			frames = append(frames, fr)
		}
		return frames, nil
	}
	mirror := func(in []*image.RGBA) []*image.RGBA {
		out := make([]*image.RGBA, len(in))
		for i, fr := range in {
			m := image.NewRGBA(fr.Bounds())
			for y := 0; y < 32; y++ {
				for x := 0; x < 32; x++ {
					m.Set(31-x, y, fr.At(x, y))
				}
			}
			out[i] = m
		}
		return out
	}
	s := &Skin{Sprites: map[string][]*image.RGBA{}}
	load := map[string]string{"protagonist": "", "protagonist-left": "", "villain": "", "villain2": "",
		"heart": "", "diamond": "", "checkpoint": "", "orb": "", "spikes": "", "arrow": "", "ladder": "", "floor": "", "wall": ""}
	got := map[string][]*image.RGBA{}
	for n := range load {
		fr, err := sheet(n)
		if err != nil {
			return nil, err
		}
		got[n] = fr
	}
	s.Sprites["sprite_right"] = got["protagonist"]
	s.Sprites["sprite_left"] = got["protagonist-left"]
	s.Sprites["sprite_climbing"] = got["protagonist"][:1]
	s.Sprites["sprite_monsterr"] = got["villain"]
	s.Sprites["sprite_monsterl"] = mirror(got["villain"])
	s.Sprites["sprite_flyingr"] = got["villain2"]
	s.Sprites["sprite_flyingl"] = mirror(got["villain2"])
	// The star coin (stelo_poentoj, +10): the twinkling star from the top of the
	// checkpoint's pole (its four frames shine orange, yellow, white, yellow),
	// centred in the cell — the game is Stelo, a star; the orb on a pole read as
	// a switch, not a thing to pick up (Joop 2026-11-02).
	s.Sprites["sprite_stelo_poentoj"] = starCoin(got["checkpoint"])
	s.Sprites["sprite_bonus"] = got["heart"]
	s.Sprites["diamanto"] = got["diamond"]
	// The checkpoint: the orb on a pole (a switch) until it is reached, the
	// shining star after (Joop 2026-11-02: the star in both states made the
	// change too subtle to see).
	s.Sprites["transirejo_na"] = got["orb"]
	s.Sprites["transirejo_a"] = got["checkpoint"]
	s.Sprites["sprite_death"] = got["spikes"]
	s.Sprites["sprite_level"] = got["arrow"]
	s.Sprites["sprite_ladder"] = got["ladder"]
	s.Sprites["sprite_block"] = got["wall"]
	// No art was made for these two; derived from the set: the flattened
	// monster is the villain squashed into the bottom 8 px (where the
	// original's sprite_monsterflat sits, box rows 27–31), the life icon the
	// protagonist shrunk to 10×10 (sprite_life's size).
	s.Sprites["sprite_monsterflat"] = []*image.RGBA{scale(got["villain"][0], 32, 8, 0, 24, 32, 32)}
	s.Sprites["sprite_life"] = []*image.RGBA{scale(got["protagonist"][0], 10, 10, 0, 0, 10, 10)}
	s.BlockTop = got["floor"]
	s.Rate = map[string]float64{
		"sprite_right": 0.2, "sprite_left": 0.2, // a walk of 6 frames a second
		"sprite_monsterr": 0.15, "sprite_monsterl": 0.15,
		"sprite_flyingr": 0.12, "sprite_flyingl": 0.12,
		"sprite_stelo_poentoj": 0.1, "transirejo_a": 0.1, // the stars twinkle
		"transirejo_na": 0.1,                   // the switch sways
		"sprite_bonus":  0.08, "diamanto": 0.1, // the heart's glow, the diamond's glint
	}
	s.IdleFirst = map[string]bool{"sprite_right": true, "sprite_left": true}
	return s, nil
}

// CoinFill and CoinRim colour the star coin.
var (
	CoinFill = color.RGBA{0x2b, 0x4c, 0x9a, 0xff}
	CoinRim  = color.RGBA{0x14, 0x24, 0x52, 0xff}
)

// defaultRate is the skin's animation speed for sprites Rate doesn't name.
const defaultRate = 0.15

// Frame is the frame to draw of an instance's skin sprite, at the skin's own
// pace from a step counter (tick), each instance offset by where it started so
// a row of coins doesn't flash together; an IdleFirst sprite standing still
// shows its first frame.
func (s *Skin) Frame(sprite string, frames []*image.RGBA, tick int, xStart, yStart, x, y, xPrev, yPrev float64) *image.RGBA {
	if len(frames) == 1 {
		return frames[0]
	}
	if s.IdleFirst[sprite] && x == xPrev && y == yPrev {
		return frames[0]
	}
	rate, ok := s.Rate[sprite]
	if !ok {
		rate = defaultRate
	}
	phase := int(xStart)/32*7 + int(yStart)/32*13
	return frames[(int(float64(tick)*rate)+phase)%len(frames)]
}

// starCoin is the star coin: the star from the top of the checkpoint's pole
// (rows 1–13, columns 10–22; the pole's grey pixels between its legs left
// out), its frames' shine kept, on a round coin of CoinFill with a CoinRim
// edge, centred in a 32×32 cell (Joop 2026-11-02: "star-marked coins").
func starCoin(checkpoint []*image.RGBA) []*image.RGBA {
	out := make([]*image.RGBA, len(checkpoint))
	for i, fr := range checkpoint {
		c := image.NewRGBA(image.Rect(0, 0, 32, 32))
		for y := 0; y < 32; y++ { // the coin: radius 10.5 around the cell's centre
			for x := 0; x < 32; x++ {
				dx, dy := float64(x)-15.5, float64(y)-15.5
				switch d := dx*dx + dy*dy; {
				case d <= 9.5*9.5:
					c.Set(x, y, CoinFill)
				case d <= 10.5*10.5:
					c.Set(x, y, CoinRim)
				}
			}
		}
		for y := 1; y <= 13; y++ {
			for x := 10; x <= 22; x++ {
				p := fr.RGBAAt(x, y)
				if p.A == 0 || isPoleGrey(p) {
					continue
				}
				c.SetRGBA(x-10+9, y-1+9, p)
			}
		}
		out[i] = c
	}
	return out
}

// isPoleGrey reports the checkpoint pole's blue-greys (the star itself is
// orange and yellow, white at its brightest).
func isPoleGrey(p color.RGBA) bool {
	return p.B > p.R && p.B > p.G
}

// scale draws src (nearest neighbour) into the w×h area at (x, y) of a new
// outW×outH image.
func scale(src *image.RGBA, w, h, x, y, outW, outH int) *image.RGBA {
	b := src.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, outW, outH))
	for yy := 0; yy < h; yy++ {
		for xx := 0; xx < w; xx++ {
			out.Set(x+xx, y+yy, src.At(b.Min.X+xx*b.Dx()/w, b.Min.Y+yy*b.Dy()/h))
		}
	}
	return out
}
