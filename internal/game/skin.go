package game

import (
	"image"
	"image/draw"
	"image/png"
	"io/fs"
)

// Skin is a drawing-only replacement of the original sprites (LOOP item 11):
// frames by original sprite name. Collision boxes stay the original's, so
// the rules don't change. Sprites it doesn't name keep the original art.
type Skin struct {
	Sprites map[string][]*image.RGBA
	// BlockTop is drawn for a block with no block right above it (a floor
	// with grass), the plain block frames elsewhere.
	BlockTop []*image.RGBA
}

// LoadNovaSkin builds the skin from Joop's commissioned 32×32 set
// (art/nova/). The roles, ours to choose (docs/DESIGN.md "New art"):
// the protagonist star is the player (its left sheet facing left; climbing
// shows its first frame), the orange villain the walking monster and the
// bonus level's (mirrored to face left), the green villain the flyer (it
// can't be stomped, so it looks meaner), the orb on a pole a star coin, the
// heart and diamond themselves, the checkpoint unused (first frame) and
// used (shining), the spikes the fire, the arrow sign the level exit, the
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
	s.Sprites["sprite_stelo_poentoj"] = got["orb"]
	s.Sprites["sprite_bonus"] = got["heart"]
	s.Sprites["diamanto"] = got["diamond"]
	s.Sprites["transirejo_na"] = got["checkpoint"][:1]
	s.Sprites["transirejo_a"] = got["checkpoint"]
	s.Sprites["sprite_death"] = got["spikes"]
	s.Sprites["sprite_level"] = got["arrow"]
	s.Sprites["sprite_ladder"] = got["ladder"]
	s.Sprites["sprite_block"] = got["wall"]
	s.BlockTop = got["floor"]
	return s, nil
}
