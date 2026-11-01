// Package level reads and writes Stelo's level files (docs/LEVEL-FORMAT.md).
package level

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Instance is one object placed in a room, as in the original (object name,
// pixel position, creation code).
type Instance struct {
	Object string
	X, Y   int
	Code   string
}

// View is a room view (Game Maker views): the visible rectangle, the object
// it follows, how close that object may come to the view's edge before it
// scrolls (BorderH/V) and the most it scrolls per step (SpeedH/V; -1 = no
// limit).
type View struct {
	X, Y, W, H       int
	Follow           string
	BorderH, BorderV int
	SpeedH, SpeedV   int
}

// Level is one room.
type Level struct {
	Name       string
	W, H       int
	Speed      int
	Persistent bool
	Colour     uint32 // RGB, 0xRRGGBB
	Background string // background image, tiled; "" = none
	View       *View
	Code       string // room creation code
	Caption    string
	Instances  []Instance
}

// GridChars maps the grid's characters to objects. Ladders have their own
// layer ('H'), so a ladder can share a cell with anything else.
var GridChars = map[byte]string{
	'#': "block", '@': "character", 'M': "monster", 'F': "flyer", 'X': "death", 'E': "levelexit",
	'*': "stelo_poentoj", '+': "bonus", '|': "marker", 'T': "transirejo", 'b': "b_monster",
	'D': "iru_al_bonuso", 'S': "monsterfarilo", 's': "finmonsterfarilo",
}

var objectChar = func() map[string]byte {
	m := map[string]byte{}
	for c, o := range GridChars {
		m[o] = c
	}
	return m
}()

const cell = 32

// Format writes the level in the text format.
func Format(l *Level) string {
	var b strings.Builder
	fmt.Fprintf(&b, "stelo-level 1\nname %s\nsize %d %d\nspeed %d\n", l.Name, l.W, l.H, l.Speed)
	fmt.Fprintf(&b, "colour #%06X\n", l.Colour)
	if l.Persistent {
		b.WriteString("persistent\n")
	}
	if l.Background != "" {
		fmt.Fprintf(&b, "background %s\n", l.Background)
	}
	if l.Caption != "" {
		fmt.Fprintf(&b, "caption %s\n", strconv.Quote(l.Caption))
	}
	if v := l.View; v != nil {
		fmt.Fprintf(&b, "view %d %d %d %d", v.X, v.Y, v.W, v.H)
		if v.Follow != "" {
			fmt.Fprintf(&b, " follow %s border %d %d speed %d %d", v.Follow, v.BorderH, v.BorderV, v.SpeedH, v.SpeedV)
		}
		b.WriteString("\n")
	}
	if l.Code != "" {
		fmt.Fprintf(&b, "code %s\n", strconv.Quote(l.Code))
	}
	gw, gh := (l.W+cell-1)/cell, (l.H+cell-1)/cell
	grid, ladders := blank(gw, gh), blank(gw, gh)
	var rest []Instance
	for _, in := range l.Instances {
		cx, cy := in.X/cell, in.Y/cell
		onGrid := in.X%cell == 0 && in.Y%cell == 0 && in.Code == "" && cx < gw && cy < gh && in.X >= 0 && in.Y >= 0
		switch c, ok := objectChar[in.Object]; {
		case onGrid && in.Object == "ladder" && ladders[cy][cx] == '.':
			ladders[cy][cx] = 'H'
		case onGrid && ok && grid[cy][cx] == '.':
			grid[cy][cx] = c
		default:
			rest = append(rest, in)
		}
	}
	b.WriteString("grid\n")
	for _, r := range grid {
		b.Write(r)
		b.WriteString("\n")
	}
	b.WriteString("ladders\n")
	for _, r := range ladders {
		b.Write(r)
		b.WriteString("\n")
	}
	for _, in := range rest {
		if t, ok := textOf(in); ok {
			fmt.Fprintf(&b, "text %d %d %s\n", in.X, in.Y, strconv.Quote(t))
			continue
		}
		fmt.Fprintf(&b, "instance %s %d %d", in.Object, in.X, in.Y)
		if in.Code != "" {
			fmt.Fprintf(&b, " %s", strconv.Quote(in.Code))
		}
		b.WriteString("\n")
	}
	b.WriteString("end\n")
	return b.String()
}

// textOf recognises a tekstilo whose creation code only sets its text.
func textOf(in Instance) (string, bool) {
	if in.Object != "tekstilo" || !strings.HasPrefix(in.Code, "teksto='") || !strings.HasSuffix(in.Code, "';") {
		return "", false
	}
	t := in.Code[len("teksto='") : len(in.Code)-2]
	if strings.Contains(t, "'") {
		return "", false
	}
	return t, true
}

func blank(w, h int) [][]byte {
	g := make([][]byte, h)
	for i := range g {
		g[i] = []byte(strings.Repeat(".", w))
	}
	return g
}

// Parse reads a level file.
func Parse(src string) (*Level, error) {
	l := &Level{}
	sc := bufio.NewScanner(strings.NewReader(src))
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	line := 0
	next := func() (string, bool) {
		if !sc.Scan() {
			return "", false
		}
		line++
		return sc.Text(), true
	}
	fail := func(f string, a ...any) error { return fmt.Errorf("line %d: %s", line, fmt.Sprintf(f, a...)) }
	if s, _ := next(); s != "stelo-level 1" {
		return nil, fail("not a stelo-level 1 file")
	}
	readGrid := func(onCell func(c byte, x, y int) error) error {
		for y := 0; y < (l.H+cell-1)/cell; y++ {
			s, ok := next()
			if !ok || len(s) != (l.W+cell-1)/cell {
				return fail("grid row %d: want %d cells", y, (l.W+cell-1)/cell)
			}
			for x := 0; x < len(s); x++ {
				if s[x] == '.' {
					continue
				}
				if err := onCell(s[x], x*cell, y*cell); err != nil {
					return err
				}
			}
		}
		return nil
	}
	for {
		s, ok := next()
		if !ok {
			return nil, fail("missing end")
		}
		key, rest, _ := strings.Cut(s, " ")
		var err error
		switch key {
		case "end":
			return l, nil
		case "name":
			l.Name = rest
		case "size":
			_, err = fmt.Sscan(rest, &l.W, &l.H)
		case "speed":
			l.Speed, err = strconv.Atoi(rest)
		case "colour":
			var c uint64
			c, err = strconv.ParseUint(strings.TrimPrefix(rest, "#"), 16, 32)
			l.Colour = uint32(c)
		case "persistent":
			l.Persistent = true
		case "background":
			l.Background = rest
		case "caption":
			l.Caption, err = strconv.Unquote(rest)
		case "view":
			v := &View{}
			f := strings.Fields(rest)
			if len(f) != 4 && len(f) != 12 {
				return nil, fail("view: x y w h [follow OBJ border H V speed H V]")
			}
			_, err = fmt.Sscan(strings.Join(f[:4], " "), &v.X, &v.Y, &v.W, &v.H)
			if err == nil && len(f) == 12 {
				if f[4] != "follow" || f[6] != "border" || f[9] != "speed" {
					return nil, fail("view: x y w h [follow OBJ border H V speed H V]")
				}
				v.Follow = f[5]
				_, err = fmt.Sscan(strings.Join([]string{f[7], f[8], f[10], f[11]}, " "), &v.BorderH, &v.BorderV, &v.SpeedH, &v.SpeedV)
			}
			l.View = v
		case "code":
			l.Code, err = strconv.Unquote(rest)
		case "grid":
			err = readGrid(func(c byte, x, y int) error {
				o, ok := GridChars[c]
				if !ok {
					return fail("unknown grid character %q", c)
				}
				l.Instances = append(l.Instances, Instance{Object: o, X: x, Y: y})
				return nil
			})
		case "ladders":
			err = readGrid(func(c byte, x, y int) error {
				if c != 'H' {
					return fail("ladders layer: only 'H'")
				}
				l.Instances = append(l.Instances, Instance{Object: "ladder", X: x, Y: y})
				return nil
			})
		case "text", "instance":
			in := Instance{Object: "tekstilo"}
			if key == "instance" {
				in.Object, rest, _ = strings.Cut(rest, " ")
			}
			f := strings.SplitN(rest, " ", 3)
			if len(f) < 2 {
				return nil, fail("%s: position missing", key)
			}
			if in.X, err = strconv.Atoi(f[0]); err == nil {
				in.Y, err = strconv.Atoi(f[1])
			}
			if err == nil && len(f) == 3 {
				var q string
				if q, err = strconv.Unquote(f[2]); err == nil {
					if key == "text" {
						in.Code = "teksto='" + q + "';"
					} else {
						in.Code = q
					}
				}
			} else if err == nil && key == "text" {
				err = fmt.Errorf("text without a string")
			}
			l.Instances = append(l.Instances, in)
		default:
			return nil, fail("unknown key %q", key)
		}
		if err != nil {
			return nil, fail("%s: %v", key, err)
		}
	}
}

// SortInstances orders instances by object, y, x, code — for comparing levels
// (the file doesn't keep the original creation order).
func SortInstances(in []Instance) {
	sort.Slice(in, func(i, j int) bool {
		a, b := in[i], in[j]
		if a.Object != b.Object {
			return a.Object < b.Object
		}
		if a.Y != b.Y {
			return a.Y < b.Y
		}
		if a.X != b.X {
			return a.X < b.X
		}
		return a.Code < b.Code
	})
}
