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
	// Cols > 0 lays the items out as blocks (Joop 2026-11-04: "square /
	// rectangle block buttons"): Cols across, each CellW × CellH, Gap apart,
	// the grid centred on X with its top at Y. Arrows move between blocks; a
	// block's Value is a second, smaller line; Enter or a click takes it.
	Cols              int
	CellW, CellH, Gap float64
}

// MenuItem is a menu line.
type MenuItem struct {
	// Action: a valued item that Enter or a click chooses (← → still change
	// its value) — Continue with its level pick.
	Action   bool
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
	if m.Cols > 0 {
		return m.updateGrid(in, res)
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
			if m.Items[line].Value != "" && !m.Items[line].Action {
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
		if m.Items[m.Sel].Value != "" && !m.Items[m.Sel].Action {
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

// updateGrid is Update for blocks: ← → step along, ↑ ↓ a row; the pointer
// selects, a click or Enter takes.
func (m *Menu) updateGrid(in MenuInput, res MenuResult) MenuResult {
	n := len(m.Items)
	step := func(d int) {
		if j := m.Sel + d; j >= 0 && j < n && !m.Items[j].Disabled {
			m.Sel = j
		}
	}
	if in.Left {
		step(-1)
	}
	if in.Right {
		step(1)
	}
	if in.Up {
		step(-m.Cols)
	}
	if in.Down {
		if m.Sel+m.Cols < n {
			step(m.Cols)
		} else if (m.Sel/m.Cols+1)*m.Cols < n { // a shorter last row: its last block
			step(n - 1 - m.Sel)
		}
	}
	if i := m.cellAt(in.PX, in.PY); i >= 0 && (in.Moved || in.Clicked) && !m.Items[i].Disabled {
		m.Sel = i
		if in.Clicked {
			res.Chosen = i
			return res
		}
	}
	if in.Enter {
		res.Chosen = m.Sel
	}
	return res
}

// cell is block i's rectangle (frame coordinates).
func (m *Menu) cell(i int) (x, y float64) {
	cols := min(m.Cols, len(m.Items))
	total := float64(cols)*m.CellW + float64(cols-1)*m.Gap
	row, col := i/m.Cols, i%m.Cols
	if last := len(m.Items) / m.Cols; row == last && len(m.Items)%m.Cols != 0 { // the last row centred
		k := len(m.Items) % m.Cols
		total = float64(k)*m.CellW + float64(k-1)*m.Gap
	}
	return m.X - total/2 + float64(col)*(m.CellW+m.Gap), m.Y + float64(row)*(m.CellH+m.Gap)
}

func (m *Menu) cellAt(px, py int) int {
	fx, fy := float64(px), float64(py)
	for i := range m.Items {
		x, y := m.cell(i)
		if fx >= x && fx < x+m.CellW && fy >= y && fy < y+m.CellH {
			return i
		}
	}
	return -1
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

// MenuColors are the menu's colours; Block the blocks' fill (grid mode).
type MenuColors struct {
	Text, Disabled, Selected, Highlight, Block color.Color
}

// drawGrid draws the blocks: each filled, the chosen one highlighted, the
// label centred and its value under it.
func (m *Menu) drawGrid(o *Overlay, c MenuColors) {
	block := c.Block
	if block == nil {
		block = color.RGBA{0, 0, 0, 120}
	}
	for i, it := range m.Items {
		x, y := m.cell(i)
		fill, col := block, c.Text
		if i == m.Sel {
			fill, col = c.Highlight, c.Selected
		}
		if it.Disabled {
			col = c.Disabled
		}
		o.AddRect(Rect{X: x, Y: y, W: m.CellW, H: m.CellH, Color: fill})
		ty := y + (m.CellH-m.Size*1.36)/2
		if it.Value != "" {
			ty -= m.Size * 0.55
			o.Add(Text{S: it.Value, X: x + m.CellW/2, Y: ty + m.Size*1.3, Size: m.Size * 0.72, Weight: Bold, Color: col, Align: Centre})
		}
		o.Add(Text{S: it.Label, X: x + m.CellW/2, Y: ty, Size: m.Size, Weight: Black, Color: col, Align: Centre})
	}
}

// Draw puts the menu into the overlay: the highlight under the chosen line,
// each line centred.
func (m *Menu) Draw(o *Overlay, c MenuColors) {
	if m.Cols > 0 {
		m.drawGrid(o, c)
		return
	}
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
