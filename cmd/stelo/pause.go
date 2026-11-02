package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/LaPingvino/stelo/internal/game"
)

// The pause menu (docs/UI.md §3, step U2b): Esc or P during a level, or the
// pause button on phones. The game does not step while it is open; the frame
// stays, dimmed. (In the 2006 original Esc ended the program.)
type pauseMenu struct {
	on     bool
	menu   game.Menu
	ids    []string
	px, py int
}

// pauseWanted: Esc or P pressed during a level.
func (a *app) pauseWanted() bool {
	if a.game == nil || a.game.Level() == 0 {
		return false
	}
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) || a.touch.pauseTapped
}

func (a *app) openPause() {
	a.pause.on = true
	a.buildPauseMenu(0)
}

func (a *app) buildPauseMenu(sel int) {
	snd := "ŝaltita"
	if soundMuted() {
		snd = "malŝaltita"
	}
	items := []game.MenuItem{
		{Label: game.Tr("Daŭrigi")},
		{Label: game.Tr("Rekomenci la nivelon")},
		{Label: game.Tr("Sono"), Value: game.Tr(snd)},
		{Label: game.Tr("Ĉefmenuo")},
	}
	a.pause.ids = []string{"resume", "restart", "sound", "title"}
	a.pause.menu = game.Menu{Items: items, Sel: min(sel, len(items)-1), X: game.ScreenW / 2, Y: 190, LineH: 30, W: 330, Size: 17}
}

// updatePause runs the open pause menu.
func (a *app) updatePause() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) || a.touch.pauseTapped {
		a.pause.on = false
		return
	}
	in := game.MenuInput{
		Up: inpututil.IsKeyJustPressed(ebiten.KeyArrowUp), Down: inpututil.IsKeyJustPressed(ebiten.KeyArrowDown),
		Left: inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft), Right: inpututil.IsKeyJustPressed(ebiten.KeyArrowRight),
		Enter: inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace),
	}
	x, y := toFrame(ebiten.CursorPosition())
	in.Moved = x != a.pause.px || y != a.pause.py
	a.pause.px, a.pause.py = x, y
	in.PX, in.PY = x, y
	in.Clicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
		in.PX, in.PY = toFrame(ebiten.TouchPosition(ids[0]))
		in.Clicked = true
	}
	res := a.pause.menu.Update(in)
	sel := a.pause.menu.Sel
	if res.Step != 0 && a.pause.ids[sel] == "sound" {
		toggleSound()
		a.buildPauseMenu(sel)
		return
	}
	if res.Chosen < 0 {
		return
	}
	switch a.pause.ids[res.Chosen] {
	case "resume":
		a.pause.on = false
	case "restart":
		a.pause.on = false
		if err := a.game.RestartLevel(); err != nil {
			log.Print(err)
		}
	case "sound":
		toggleSound()
		a.buildPauseMenu(sel)
	case "title":
		a.pause.on = false
		a.sound.StopMusic()
		a.game = nil
		a.openTitle()
	}
}

// drawPause dims the frame and puts the menu over it.
func (a *app) drawPause() {
	o := &a.overlay
	o.AddRect(game.Rect{X: 0, Y: 0, W: game.ScreenW, H: game.ScreenH, Color: color.RGBA{0, 0, 0, 140}})
	o.AddRect(game.Rect{X: game.ScreenW/2 - 200, Y: 72, W: 400, H: 250, Color: color.RGBA{0x2a, 0x1c, 0x12, 225}}) // a panel behind the menu
	o.Add(game.Text{S: game.Tr("Paŭzo"), X: game.ScreenW / 2, Y: 110, Size: 44, Weight: game.Black, Color: color.RGBA{0xff, 0xd0, 0x3a, 255}, Align: game.Centre})
	a.pause.menu.Draw(o, game.MenuColors{Text: color.RGBA{0xff, 0xf4, 0xd8, 255}, Disabled: color.RGBA{0xb0, 0xa0, 0x80, 255},
		Selected: color.RGBA{0x2a, 0x1c, 0x12, 255}, Highlight: color.RGBA{0xff, 0xd0, 0x3a, 230}})
}
