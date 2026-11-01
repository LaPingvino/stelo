package level

import (
	"encoding/json"
	"fmt"
	"os"
)

// The parts of original/project.json (gm53-reader) the converter reads.
type gmProject struct {
	Objects     []*struct{ Name string }
	Backgrounds []*struct{ Name string }
	Rooms       []*struct {
		Name, Caption string
		CreationCode  string `json:"creation_code"`
		W, H, Speed   int
		Persistent    bool
		BgColor       uint32 `json:"bg_color"`
		Backgrounds   []struct {
			Visible bool
			TileH   bool `json:"tile_h"`
			TileV   bool `json:"tile_v"`
			Bg      int
		}
		ViewsEnabled bool `json:"views_enabled"`
		Views        []struct {
			Visible bool
			View    [4]int
			Border  [4]int `json:"border_hv_speed_hv"`
			Follow  int
		}
		Instances []struct {
			X, Y, Obj int
			Code      string
		}
	}
}

// LoadOriginal reads every room of the recovered project as a Level.
func LoadOriginal(projectJSON string) ([]*Level, error) {
	data, err := os.ReadFile(projectJSON)
	if err != nil {
		return nil, err
	}
	var p gmProject
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	obj := func(i int) string {
		if i >= 0 && i < len(p.Objects) && p.Objects[i] != nil {
			return p.Objects[i].Name
		}
		return fmt.Sprintf("object%d", i)
	}
	var out []*Level
	for _, r := range p.Rooms {
		if r == nil {
			continue
		}
		c := r.BgColor // Game Maker colours are 0xBBGGRR
		l := &Level{Name: r.Name, W: r.W, H: r.H, Speed: r.Speed, Persistent: r.Persistent,
			Colour: (c&0xFF)<<16 | (c & 0xFF00) | (c>>16)&0xFF, Code: r.CreationCode, Caption: r.Caption}
		for _, b := range r.Backgrounds {
			if b.Visible && b.Bg >= 0 && b.Bg < len(p.Backgrounds) {
				if !b.TileH || !b.TileV {
					return nil, fmt.Errorf("%s: background not tiled both ways (not in the format yet)", r.Name)
				}
				l.Background = p.Backgrounds[b.Bg].Name
				break
			}
		}
		if r.ViewsEnabled {
			for _, v := range r.Views {
				if v.Visible {
					l.View = &View{X: v.View[0], Y: v.View[1], W: v.View[2], H: v.View[3]}
					if v.Follow >= 0 {
						l.View.Follow = obj(v.Follow)
						l.View.BorderH, l.View.BorderV, l.View.SpeedH, l.View.SpeedV = v.Border[0], v.Border[1], v.Border[2], v.Border[3]
					}
					break
				}
			}
		}
		for _, in := range r.Instances {
			l.Instances = append(l.Instances, Instance{Object: obj(in.Obj), X: in.X, Y: in.Y, Code: in.Code})
		}
		out = append(out, l)
	}
	return out, nil
}
