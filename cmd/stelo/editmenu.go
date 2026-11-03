package main

import (
	"fmt"
	"image/color"
	"log"
	"runtime"
	"slices"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/LaPingvino/stelo/internal/editor"
	"github.com/LaPingvino/stelo/internal/game"
	"github.com/LaPingvino/stelo/internal/level"
	"github.com/LaPingvino/stelo/internal/store"
)

// The editor's menu (LOOP item UE; Joop 2026-11-03: "the editor doesn't seem
// to have an exit or existing level picker"): Esc or the strip's ☰ button.
// Nova nivelo · Malfermi ‹level› (the original levels, the bonus level, the
// levels saved here) · El dosiero… · Konservi · Ludi · Ĉefmenuo.
type editMenu struct {
	on     bool
	menu   game.Menu
	ids    []string
	levels []string // what Malfermi steps through
	pick   int
	px, py int
}

// editLevels: the original levels, then the saved ones that aren't one of them.
func editLevels() []string {
	names := []string{"nivelo1", "nivelo2", "nivelo3", "nivelo4", "nivelo5", "bonusnivelo"}
	for _, n := range store.SavedLevels() {
		if !slices.Contains(names, n) {
			names = append(names, n)
		}
	}
	return names
}

func (a *app) openEditMenu() {
	m := &a.edit.menu
	m.on = true
	m.levels = editLevels()
	m.pick = max(slices.Index(m.levels, a.edit.screen.Doc.L.Name), 0)
	a.buildEditMenu(0)
}

func (a *app) buildEditMenu(sel int) {
	m := &a.edit.menu
	items := []game.MenuItem{
		{Label: game.Tr("Nova nivelo")},
		{Label: game.Tr("Malfermi"), Value: m.levels[m.pick]},
	}
	m.ids = []string{"new", "open"}
	if runtime.GOOS == "js" { // the browser's file picker (the desktop drops files on the window)
		items = append(items, game.MenuItem{Label: game.Tr("El dosiero…")})
		m.ids = append(m.ids, "file")
	}
	items = append(items, game.MenuItem{Label: game.Tr("Konservi")}, game.MenuItem{Label: game.Tr("Ludi")},
		game.MenuItem{Label: game.Tr("Daŭrigi redakti")}, game.MenuItem{Label: game.Tr("Ĉefmenuo")})
	m.ids = append(m.ids, "save", "play", "back", "title")
	m.menu = game.Menu{Items: items, Sel: min(sel, len(items)-1), X: game.ScreenW / 2, Y: 130, LineH: 30, W: 360, Size: 17}
}

func (a *app) updateEditMenu() {
	m := &a.edit.menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		m.on = false
		return
	}
	in := game.MenuInput{
		Up: inpututil.IsKeyJustPressed(ebiten.KeyArrowUp), Down: inpututil.IsKeyJustPressed(ebiten.KeyArrowDown),
		Left: inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft), Right: inpututil.IsKeyJustPressed(ebiten.KeyArrowRight),
		Enter: inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace),
	}
	x, y := toFrame(ebiten.CursorPosition())
	in.Moved = x != m.px || y != m.py
	m.px, m.py = x, y
	in.PX, in.PY = x, y
	in.Clicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
		in.PX, in.PY = toFrame(ebiten.TouchPosition(ids[0]))
		in.Clicked = true
	}
	res := m.menu.Update(in)
	sel := m.menu.Sel
	if res.Step != 0 && m.ids[sel] == "open" {
		m.pick = (m.pick + res.Step + len(m.levels)) % len(m.levels)
		a.buildEditMenu(sel)
		return
	}
	if res.Chosen < 0 {
		return
	}
	m.on = false
	switch m.ids[res.Chosen] {
	case "new":
		a.newLevel()
	case "open":
		a.openEditor(m.levels[m.pick])
	case "file":
		a.pickLevelFile()
	case "save":
		a.saveEdited()
	case "play":
		a.playEdited()
	case "title":
		a.edit.on = false
		a.game = nil
		a.openTitle()
	}
}

func (a *app) drawEditMenu() {
	o := &a.overlay
	o.AddRect(game.Rect{X: 0, Y: 0, W: game.ScreenW, H: game.ScreenH, Color: color.RGBA{0, 0, 0, 120}})
	o.AddRect(game.Rect{X: game.ScreenW/2 - 210, Y: 60, W: 420, H: 320, Color: color.RGBA{0x2a, 0x1c, 0x12, 230}})
	o.Add(game.Text{S: game.Tr("Redaktilo"), X: game.ScreenW / 2, Y: 76, Size: 30, Weight: game.Black, Color: color.RGBA{0xff, 0xd0, 0x3a, 255}, Align: game.Centre})
	m := &a.edit.menu
	m.menu.Draw(o, game.MenuColors{Text: color.RGBA{0xff, 0xf4, 0xd8, 255}, Disabled: color.RGBA{0xb0, 0xa0, 0x80, 255},
		Selected: color.RGBA{0x2a, 0x1c, 0x12, 255}, Highlight: color.RGBA{0xff, 0xd0, 0x3a, 230}})
}

// newLevelText is an empty level: walls round the edge, the player bottom
// left, the exit bottom right, the life controller (so the HUD shows).
func newLevelText(name string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "stelo-level 1\nname %s\nsize 640 480\nspeed 30\ncolour #00FF80\nbackground background_sky\n", name)
	b.WriteString("view 0 0 640 480 follow character border 64 64 speed 4 4\ngrid\n")
	for y := range 15 {
		switch y {
		case 0, 14:
			b.WriteString(strings.Repeat("#", 20))
		case 13:
			b.WriteString("#@" + strings.Repeat(".", 16) + "E#")
		default:
			b.WriteString("#" + strings.Repeat(".", 18) + "#")
		}
		b.WriteByte('\n')
	}
	b.WriteString("instance life_controller 32 0\nend\n")
	return b.String()
}

// newLevel starts a new level, named nova1, nova2, … (the first free name).
func (a *app) newLevel() {
	taken := editLevels()
	name := ""
	for i := 1; ; i++ {
		name = fmt.Sprintf("nova%d", i)
		if !slices.Contains(taken, name) && a.edit.docs[name] == nil {
			break
		}
	}
	l, err := level.Parse(newLevelText(name))
	if err != nil {
		log.Printf("editor: new level: %v", err)
		return
	}
	a.edit.docs[name] = editor.New(l)
	a.openEditor(name)
}
