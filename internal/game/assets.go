package game

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/fs"
	"path"
)

// Sprite is one sprite of the original: its frames as RGBA (the transparent
// colour already cleared) and its bounding box.
type Sprite struct {
	Name   string
	Frames []*image.RGBA
	// BBox is left, right, top, bottom inside the 32×32 cell (inclusive).
	BBox [4]int
}

// ObjectDef is what the game needs of an original object.
type ObjectDef struct {
	Name    string
	Sprite  string
	Visible bool
	Depth   int
	Solid   bool
}

// Assets are the loaded sprites, backgrounds and object definitions.
type Assets struct {
	Sprites     map[string]*Sprite
	Backgrounds map[string]*image.RGBA
	Objects     map[string]*ObjectDef
}

type project struct {
	Sprites []*struct {
		Name        string
		Transparent bool
		Bbox        [4]int `json:"bbox_lrbt"`
		Frames      []string
	}
	Backgrounds []*struct {
		Name        string
		Transparent bool
		File        string
	}
	Objects []*struct {
		Name    string
		Sprite  int
		Visible bool
		Depth   int
		Solid   bool
	}
}

// LoadAssets reads the original's resources from fsys (the stelo.Files layout).
func LoadAssets(fsys fs.FS) (*Assets, error) {
	data, err := fs.ReadFile(fsys, "original/project.json")
	if err != nil {
		return nil, err
	}
	var p project
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	a := &Assets{Sprites: map[string]*Sprite{}, Backgrounds: map[string]*image.RGBA{}, Objects: map[string]*ObjectDef{}}
	for _, s := range p.Sprites {
		if s == nil {
			continue
		}
		// bbox_lrbt is stored left, right, bottom, top.
		sp := &Sprite{Name: s.Name, BBox: [4]int{s.Bbox[0], s.Bbox[1], s.Bbox[3], s.Bbox[2]}}
		for _, f := range s.Frames {
			img, err := loadPNG(fsys, path.Join("original", f), s.Transparent)
			if err != nil {
				return nil, fmt.Errorf("sprite %s: %w", s.Name, err)
			}
			sp.Frames = append(sp.Frames, img)
		}
		a.Sprites[s.Name] = sp
	}
	for _, b := range p.Backgrounds {
		if b == nil || b.File == "" {
			continue
		}
		img, err := loadPNG(fsys, path.Join("original", b.File), b.Transparent)
		if err != nil {
			return nil, fmt.Errorf("background %s: %w", b.Name, err)
		}
		a.Backgrounds[b.Name] = img
	}
	for _, o := range p.Objects {
		if o == nil {
			continue
		}
		d := &ObjectDef{Name: o.Name, Visible: o.Visible, Depth: o.Depth, Solid: o.Solid}
		if o.Sprite >= 0 && o.Sprite < len(p.Sprites) && p.Sprites[o.Sprite] != nil {
			d.Sprite = p.Sprites[o.Sprite].Name
		}
		a.Objects[o.Name] = d
	}
	return a, nil
}

// loadPNG decodes an image; with transparent set, every pixel of the colour of
// the bottom-left pixel becomes transparent (Game Maker's transparent sprites).
func loadPNG(fsys fs.FS, name string, transparent bool) (*image.RGBA, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	b := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), src, b.Min, draw.Src)
	if transparent {
		key := img.RGBAAt(0, b.Dy()-1)
		for i := 0; i < len(img.Pix); i += 4 {
			if img.Pix[i] == key.R && img.Pix[i+1] == key.G && img.Pix[i+2] == key.B {
				img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = 0, 0, 0, 0
			}
		}
	}
	return img, nil
}
