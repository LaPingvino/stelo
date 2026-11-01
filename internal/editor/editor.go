// Package editor is the level editor's model (LOOP item 10, docs/DESIGN.md
// §3): a level edited cell by cell on the 32 px grid, with undo. The screens
// (in-game mode, its own web page) draw and drive it.
package editor

import (
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"github.com/LaPingvino/stelo/internal/level"
)

// Cell is the grid size.
const Cell = 32

// Tools: an object name from level.GridChars paints that object; Ladder
// toggles a ladder (its own layer: it can share a cell); Erase clears the
// cell.
const (
	Ladder = "ladder"
	Erase  = "erase"
)

// Unique objects: painting one moves it (a level has one player).
var unique = map[string]bool{"character": true}

// Doc is a level being edited.
type Doc struct {
	L    *level.Level
	undo [][]level.Instance
	size [][2]int // the size before each undo step
}

// New edits l (a copy is not made: the doc owns it).
func New(l *level.Level) *Doc { return &Doc{L: l} }

// onGrid reports an instance sitting exactly on cell (cx, cy).
func onGrid(in level.Instance, cx, cy int) bool {
	return in.X == cx*Cell && in.Y == cy*Cell && in.Code == ""
}

// gridObject reports whether the object has a grid character (the editor's
// palette; instances of other objects, like the life controller, are kept
// as they are).
func gridObject(name string) bool {
	for _, o := range level.GridChars {
		if o == name {
			return true
		}
	}
	return name == "ladder"
}

// At returns the grid object in the cell ("" none) and whether a ladder is
// there.
func (d *Doc) At(cx, cy int) (obj string, ladder bool) {
	for _, in := range d.L.Instances {
		if !onGrid(in, cx, cy) {
			continue
		}
		if in.Object == "ladder" {
			ladder = true
		} else if gridObject(in.Object) {
			obj = in.Object
		}
	}
	return obj, ladder
}

func (d *Doc) save() {
	d.undo = append(d.undo, slices.Clone(d.L.Instances))
	d.size = append(d.size, [2]int{d.L.W, d.L.H})
}

// inside reports a cell within the level.
func (d *Doc) inside(cx, cy int) bool {
	return cx >= 0 && cy >= 0 && cx*Cell < d.L.W && cy*Cell < d.L.H
}

// Paint applies a tool to a cell. It reports whether anything changed.
func (d *Doc) Paint(cx, cy int, tool string) bool {
	if !d.inside(cx, cy) {
		return false
	}
	obj, ladder := d.At(cx, cy)
	switch {
	case tool == Ladder:
		d.save()
		if ladder {
			d.remove(func(in level.Instance) bool { return onGrid(in, cx, cy) && in.Object == "ladder" })
		} else {
			d.L.Instances = append(d.L.Instances, level.Instance{Object: "ladder", X: cx * Cell, Y: cy * Cell})
		}
		return true
	case tool == Erase:
		if obj == "" && !ladder {
			return false
		}
		d.save()
		d.remove(func(in level.Instance) bool { return onGrid(in, cx, cy) && gridObject(in.Object) })
		return true
	case !gridObject(tool):
		return false
	case obj == tool:
		return false
	}
	d.save()
	d.remove(func(in level.Instance) bool {
		return (onGrid(in, cx, cy) && gridObject(in.Object) && in.Object != "ladder") || (unique[tool] && in.Object == tool)
	})
	d.L.Instances = append(d.L.Instances, level.Instance{Object: tool, X: cx * Cell, Y: cy * Cell})
	return true
}

func (d *Doc) remove(drop func(level.Instance) bool) {
	d.L.Instances = slices.DeleteFunc(d.L.Instances, drop)
}

// Resize sets the level's size in cells; what falls outside is dropped.
func (d *Doc) Resize(cw, ch int) bool {
	if cw < 1 || ch < 1 || (cw*Cell == d.L.W && ch*Cell == d.L.H) {
		return false
	}
	d.save()
	d.L.W, d.L.H = cw*Cell, ch*Cell
	d.remove(func(in level.Instance) bool { return in.X >= d.L.W || in.Y >= d.L.H })
	if v := d.L.View; v != nil {
		v.W, v.H = min(v.W, d.L.W), min(v.H, d.L.H)
	}
	return true
}

// Undo takes back the last change; false when there is none.
func (d *Doc) Undo() bool {
	n := len(d.undo)
	if n == 0 {
		return false
	}
	d.L.Instances = d.undo[n-1]
	d.L.W, d.L.H = d.size[n-1][0], d.size[n-1][1]
	d.undo, d.size = d.undo[:n-1], d.size[:n-1]
	return true
}

// Text is the level file (docs/LEVEL-FORMAT.md).
func (d *Doc) Text() string { return level.Format(d.L) }

// FindLevel reads the first *.txt in files that is a level file (files
// dropped on the window).
func FindLevel(files fs.FS) (*level.Level, error) {
	var found *level.Level
	var lastErr error
	_ = fs.WalkDir(files, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || found != nil || !strings.HasSuffix(strings.ToLower(p), ".txt") {
			return nil
		}
		data, err := fs.ReadFile(files, p)
		if err != nil {
			lastErr = err
			return nil
		}
		if l, err := level.Parse(string(data)); err == nil {
			found = l
		} else {
			lastErr = fmt.Errorf("%s: %w", p, err)
		}
		return nil
	})
	if found == nil {
		if lastErr == nil {
			lastErr = errors.New("no level file (*.txt)")
		}
		return nil, lastErr
	}
	return found, nil
}
