package game

import "image/color"

// Menu is a list of lines to choose from, as in Pits (docs/UI.md §3): arrows
// up / down move the choice (skipping disabled lines, wrapping), Enter takes
// it; the mouse or a finger on a line selects it and a click or tap takes it.
// A line with a Value (a setting: "Lingvo: Esperanto") changes it with left /
// right, reported as a step of −1 / +1.
type Menu struct {
	Items []MenuItem
	Sel   int
	// X, Y is the top-centre of the first line, LineH the line height, W the
	// highlight's width (frame coordinates); Size the font size.
	X, Y, LineH, W, Size float64
}

// MenuItem is a menu line.
type MenuItem struct {
	Label    string
	Value    string // shown after the label when set
	Disabled bool
}

// MenuInput is what the menu reads in a frame: key presses, and the pointer
// (frame coordinates; Moved when it moved, Clicked when it was pressed).
type MenuInput struct {
	Up, Down, Left, Right, Enter bool
	PX, PY                       int
	Moved, Clicked               bool
}

// MenuResult is what happened: Chosen the line taken (−1 none), Step a
// setting's change (−1, +1) on line Sel.
type MenuResult struct {
	Chosen int
	Step   int
}

// Update handles a frame's input.
func (m *Menu) Update(in MenuInput) MenuResult {
	res := MenuResult{Chosen: -1}
	if len(m.Items) == 0 {
		return res
	}
	if m.Items[m.Sel].Disabled {
		m.move(1)
	}
	if in.Up {
		m.move(-1)
	}
	if in.Down {
		m.move(1)
	}
	if line := m.lineAt(in.PX, in.PY); line >= 0 && (in.Moved || in.Clicked) && !m.Items[line].Disabled {
		m.Sel = line
		if in.Clicked {
			if m.Items[line].Value != "" {
				res.Step = 1 // a click on a setting steps it on
			} else {
				res.Chosen = line
			}
			return res
		}
	}
	if m.Items[m.Sel].Value != "" && (in.Left || in.Right) {
		res.Step = map[bool]int{true: -1, false: 1}[in.Left]
	}
	if in.Enter {
		if m.Items[m.Sel].Value != "" {
			res.Step = 1
		} else {
			res.Chosen = m.Sel
		}
	}
	return res
}

func (m *Menu) move(d int) {
	for range m.Items {
		m.Sel = (m.Sel + d + len(m.Items)) % len(m.Items)
		if !m.Items[m.Sel].Disabled {
			return
		}
	}
}

// lineAt is the line under (x, y), or −1.
func (m *Menu) lineAt(x, y int) int {
	fx, fy := float64(x), float64(y)
	if fx < m.X-m.W/2 || fx >= m.X+m.W/2 || fy < m.Y {
		return -1
	}
	i := int((fy - m.Y) / m.LineH)
	if i >= len(m.Items) {
		return -1
	}
	return i
}

// MenuColors are the menu's colours.
type MenuColors struct {
	Text, Disabled, Selected, Highlight color.Color
}

// Draw puts the menu into the overlay: the highlight under the chosen line,
// each line centred.
func (m *Menu) Draw(o *Overlay, c MenuColors) {
	for i, it := range m.Items {
		y := m.Y + float64(i)*m.LineH
		col, w := c.Text, Bold
		if it.Disabled {
			col, w = c.Disabled, Regular
		}
		if i == m.Sel {
			o.AddRect(Rect{X: m.X - m.W/2, Y: y, W: m.W, H: m.LineH, Color: c.Highlight})
			col = c.Selected
		}
		s := it.Label
		if it.Value != "" {
			s += ":  ‹ " + it.Value + " ›"
		}
		o.Add(Text{S: s, X: m.X, Y: y + (m.LineH-m.Size*1.36)/2, Size: m.Size, Weight: w, Color: col, Align: Centre})
	}
}
