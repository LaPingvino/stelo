package game

import (
	"image/color"
	"testing"
)

func testMenu() *Menu {
	return &Menu{Items: []MenuItem{{Label: "Ludi"}, {Label: "Daŭrigi", Disabled: true}, {Label: "Lingvo", Value: "Esperanto"}, {Label: "Fini"}},
		X: 320, Y: 200, LineH: 30, W: 300, Size: 18}
}

// Keys: down skips the disabled line and wraps; Enter takes a line; on a
// setting, left / right / Enter step it instead.
func TestMenuKeys(t *testing.T) {
	m := testMenu()
	if r := m.Update(MenuInput{Down: true}); m.Sel != 2 || r.Chosen != -1 {
		t.Fatalf("down: sel %d %+v", m.Sel, r)
	}
	if r := m.Update(MenuInput{Left: true}); r.Step != -1 {
		t.Errorf("left on a setting: %+v", r)
	}
	if r := m.Update(MenuInput{Enter: true}); r.Step != 1 || r.Chosen != -1 {
		t.Errorf("enter on a setting: %+v", r)
	}
	m.Update(MenuInput{Down: true})
	m.Update(MenuInput{Down: true}) // wraps past Fini to Ludi
	if m.Sel != 0 {
		t.Errorf("wrap: sel %d", m.Sel)
	}
	m.Update(MenuInput{Up: true}) // up from Ludi wraps to Fini
	if r := m.Update(MenuInput{Enter: true}); m.Sel != 3 || r.Chosen != 3 {
		t.Errorf("up + enter: sel %d %+v", m.Sel, r)
	}
}

// Pointer: moving over a line selects it, a click takes it; disabled lines and
// places beside the menu do nothing.
func TestMenuPointer(t *testing.T) {
	m := testMenu()
	if r := m.Update(MenuInput{PX: 320, PY: 200 + 3*30 + 5, Moved: true}); m.Sel != 3 || r.Chosen != -1 {
		t.Fatalf("hover: sel %d %+v", m.Sel, r)
	}
	if r := m.Update(MenuInput{PX: 320, PY: 200 + 30 + 5, Clicked: true}); r.Chosen != -1 || m.Sel != 3 {
		t.Errorf("click on a disabled line: sel %d %+v", m.Sel, r)
	}
	if r := m.Update(MenuInput{PX: 100, PY: 205, Clicked: true}); r.Chosen != -1 {
		t.Errorf("click beside: %+v", r)
	}
	if r := m.Update(MenuInput{PX: 320, PY: 205, Clicked: true}); r.Chosen != 0 {
		t.Errorf("click on Ludi: %+v", r)
	}
	if r := m.Update(MenuInput{PX: 320, PY: 200 + 2*30 + 5, Clicked: true}); r.Step != 1 || m.Sel != 2 {
		t.Errorf("click on a setting: %+v", r)
	}
}

// Draw: one highlight under the selected line, a text per line, the setting
// with its value.
func TestMenuDraw(t *testing.T) {
	m := testMenu()
	var o Overlay
	m.Draw(&o, MenuColors{Text: color.White, Disabled: color.Gray{128}, Selected: color.Black, Highlight: color.White})
	if len(o.Rects) != 1 || len(o.Texts) != 4 || o.Texts[2].S != "Lingvo:  ‹ Esperanto ›" {
		t.Fatalf("%d rects, texts %+v", len(o.Rects), o.Texts)
	}
}

// An Action item with a value: ← → step it, Enter chooses it.
func TestMenuActionItem(t *testing.T) {
	m := Menu{Items: []MenuItem{{Label: "Play"}, {Label: "Continue", Value: "level 3", Action: true}}, Sel: 1}
	if r := m.Update(MenuInput{Left: true}); r.Step != -1 || r.Chosen >= 0 {
		t.Errorf("left: %+v", r)
	}
	if r := m.Update(MenuInput{Enter: true}); r.Chosen != 1 || r.Step != 0 {
		t.Errorf("enter: %+v", r)
	}
}
