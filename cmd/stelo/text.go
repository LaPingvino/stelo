package main

import (
	"image"
	"io/fs"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/LaPingvino/stelo/internal/game"
)

// The overlay's texts (game.Overlay) drawn at the screen's real resolution in
// Noto Sans (fonts/, SIL Open Font License), placed by the frame's position
// on the screen (view), so they are sharp however large the window is. Uses
// Ebitengine's text package on golang.org/x/image/font (already a dependency;
// text/v2 would bring in new modules).

var noto [3]*opentype.Font

// faces caches a face per weight and pixel size.
var faces = map[[2]int]font.Face{}

func loadNoto(fsys fs.FS) {
	for i, name := range []string{"NotoSans-Regular", "NotoSans-Bold", "NotoSans-Black"} {
		data, err := fs.ReadFile(fsys, "fonts/"+name+".ttf")
		if err != nil {
			log.Printf("font %s: %v", name, err)
			continue
		}
		f, err := opentype.Parse(data)
		if err != nil {
			log.Printf("font %s: %v", name, err)
			continue
		}
		noto[i] = f
	}
}

func face(w game.Weight, px int) font.Face {
	k := [2]int{int(w), px}
	if f, ok := faces[k]; ok {
		return f
	}
	if noto[w] == nil {
		return nil
	}
	f, err := opentype.NewFace(noto[w], &opentype.FaceOptions{Size: float64(px), DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		return nil
	}
	faces[k] = f
	return f
}

// drawOverlay draws the frame's rectangles, then its texts, onto the screen.
// pics are the overlay's pictures as GPU images (made once each).
var pics = map[*image.RGBA]*ebiten.Image{}

func drawOverlay(screen *ebiten.Image, o *game.Overlay) {
	for _, r := range o.Rects {
		vector.FillRect(screen, float32(view.ox+r.X*view.scale), float32(view.oy+r.Y*view.scale),
			float32(r.W*view.scale), float32(r.H*view.scale), r.Color, false)
	}
	for _, p := range o.Pics {
		img := pics[p.Img]
		if img == nil {
			img = ebiten.NewImageFromImage(p.Img)
			pics[p.Img] = img
		}
		op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
		b := p.Img.Bounds()
		op.GeoM.Scale(p.W*view.scale/float64(b.Dx()), p.H*view.scale/float64(b.Dy()))
		op.GeoM.Translate(view.ox+p.X*view.scale, view.oy+p.Y*view.scale)
		screen.DrawImage(img, op)
	}
	for _, t := range o.Texts {
		f := face(t.Weight, int(math.Round(t.Size*view.scale)))
		if f == nil {
			continue
		}
		x := view.ox + t.X*view.scale
		y := view.oy + t.Y*view.scale + float64(f.Metrics().Ascent.Round())
		if t.Align != game.Left {
			w := float64(font.MeasureString(f, t.S).Round())
			if t.Align == game.Centre {
				x -= w / 2
			} else {
				x -= w
			}
		}
		text.Draw(screen, t.S, f, int(math.Round(x)), int(math.Round(y)), t.Color)
	}
}
