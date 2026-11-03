package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/game"
	"github.com/LaPingvino/stelo/internal/sound"
	"github.com/LaPingvino/stelo/internal/store"
)

// The title screen (docs/UI.md §3, Joop: "make it new"): the game's own art and
// a menu in Noto Sans, instead of the 2006 title picture and rules room. Play
// starts at level 1, Continue at the furthest level reached in the chosen pack
// ("resume wherever you got, like with Pits"); language and art are set here.

type titleScreen struct {
	on       bool
	menu     game.Menu
	ids      []string
	tick     int
	px, py   int // the pointer, frame coordinates
	langs    []string
	progress int    // the chosen pack's furthest level
	info     string // an open info page: "help" (Kiel ludi), "about" (Pri), "" none
	pick     int    // the level Continue starts (1 … progress; the furthest by default)
}

func (a *app) packName() string {
	if a.pack < len(a.packs) {
		return a.packs[a.pack].Name
	}
	return "Stelo"
}

// openTitle shows the title (at start, and when a game ends).
func (a *app) openTitle() {
	a.title.on = true
	a.title.langs = []string{"eo", "en"}
	a.buildTitleMenu(0)
	if a.sound == nil {
		a.sound = sound.NewPlayer(stelo.Files)
		theSound = a.sound
		a.sound.SetMuted(store.Muted())
		a.sound.SetMusicOff(store.MusicOff())
	}
	a.sound.Play("@sound:komenco") // the 2006 title tune
	ebiten.SetWindowTitle("Stelo")
}

func (a *app) buildTitleMenu(sel int) {
	a.title.progress = store.Progress(a.packName())
	var items []game.MenuItem
	var ids []string
	add := func(id, label, value string) {
		items = append(items, game.MenuItem{Label: game.Tr(label), Value: value})
		ids = append(ids, id)
	}
	add("play", "Ludi", "")
	if a.title.progress > 1 { // Continue at a reached level: the furthest, or ← → for an earlier one
		if a.title.pick < 1 || a.title.pick > a.title.progress {
			a.title.pick = a.title.progress
		}
		items = append(items, game.MenuItem{Label: game.Tr("Daŭrigi"), Value: game.Tr("nivelo") + " " + itoa(a.title.pick), Action: true})
		ids = append(ids, "continue")
	}
	if len(a.packs) > 1 {
		add("pack", "Niveloj", a.packName())
	}
	lang := "Esperanto"
	if game.Lang == "en" {
		lang = "English"
	}
	add("lang", "Lingvo", lang)
	art := game.Tr("nova")
	if a.skin == nil {
		art = game.Tr("de 2006")
	}
	add("art", "Arto", art)
	onOff := func(off bool) string {
		if off {
			return game.Tr("malŝaltita")
		}
		return game.Tr("ŝaltita")
	}
	add("sound", "Sono", onOff(soundMuted()))
	add("music", "Muziko", onOff(theSound != nil && theSound.MusicOff()))
	add("help", "Kiel ludi", "")
	add("about", "Pri", "")
	add("edit", "Redaktilo", "")
	if runtime.GOOS != "js" {
		add("quit", "Fini", "")
	}
	a.title.ids = ids
	a.title.menu = game.Menu{Items: items, Sel: min(sel, len(items)-1), X: game.ScreenW / 2, Y: 218, LineH: 26, W: 330, Size: 17}
}

func (a *app) updateTitle() {
	a.title.tick++
	in := game.MenuInput{
		Up: inpututil.IsKeyJustPressed(ebiten.KeyArrowUp), Down: inpututil.IsKeyJustPressed(ebiten.KeyArrowDown),
		Left: inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft), Right: inpututil.IsKeyJustPressed(ebiten.KeyArrowRight),
		Enter: inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace),
	}
	x, y := toFrame(ebiten.CursorPosition())
	in.Moved = x != a.title.px || y != a.title.py
	a.title.px, a.title.py = x, y
	in.PX, in.PY = x, y
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		in.Clicked = true
	}
	if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
		in.PX, in.PY = toFrame(ebiten.TouchPosition(ids[0]))
		in.Clicked = true
	}
	if a.title.info != "" { // an info page: Enter, Esc or a click goes back
		if in.Enter || in.Clicked || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			a.title.info = ""
		}
		return
	}
	res := a.title.menu.Update(in)
	sel := a.title.menu.Sel
	if res.Step != 0 {
		switch a.title.ids[sel] {
		case "continue": // the level pick: 1 … the furthest reached, round
			a.title.pick = (a.title.pick-1+res.Step+a.title.progress)%a.title.progress + 1
		case "pack":
			a.pack = (a.pack + res.Step + len(a.packs)) % len(a.packs)
			a.title.pick = 0 // the new pack's furthest level
		case "lang":
			if game.Lang == "en" {
				game.Lang = "eo"
			} else {
				game.Lang = "en"
			}
		case "art":
			a.toggleSkin()
		case "sound":
			toggleSound()
		case "music":
			toggleMusic()
		}
		a.buildTitleMenu(sel)
		return
	}
	if res.Chosen < 0 {
		return
	}
	switch a.title.ids[res.Chosen] {
	case "sound", "music": // a click or Enter changes the setting too
		if a.title.ids[res.Chosen] == "sound" {
			toggleSound()
		} else {
			toggleMusic()
		}
		a.buildTitleMenu(res.Chosen)
	case "play":
		a.startFromTitle(1)
	case "continue":
		a.startFromTitle(a.title.pick)
	case "help", "about":
		a.title.info = a.title.ids[res.Chosen]
	case "edit":
		a.title.on = false
		a.openEditor("nivelo1")
	case "quit":
		os.Exit(0)
	}
}

// startFromTitle starts the chosen pack at level n, with the game music.
func (a *app) startFromTitle(n int) {
	if a.pack == 0 {
		a.load(0)
	} else {
		a.playPack(a.packs[a.pack])
	}
	if a.game == nil {
		return
	}
	if err := a.game.EnterLevel(n); err != nil {
		log.Printf("level %d: %v", n, err)
	}
	a.game.OnRestart = func() { a.game = nil; a.openTitle() }
	a.title.on = false
	a.sound.Play("@music:muziko1")
}

// drawTitle draws the title: the frame (background, a floor, the cast) and the
// overlay (the name, a line under it, the menu).
func (a *app) drawTitle() {
	dst := a.frame
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.RGBA{0x2a, 0x1c, 0x12, 255}), image.Point{}, draw.Src)
	if bg := a.assets.Backgrounds["background_sky"]; bg != nil {
		b := bg.Bounds()
		for y := 0; y < game.ScreenH; y += b.Dy() {
			for x := 0; x < game.ScreenW; x += b.Dx() {
				draw.Draw(dst, b.Add(image.Pt(x, y)), bg, b.Min, draw.Src)
			}
		}
	}
	frame := func(name string, i int) *image.RGBA {
		if a.skin != nil {
			if fr := a.skin.Sprites[name]; len(fr) > 0 {
				return a.skin.Frame(name, fr, a.title.tick, float64(i*96), 0, 1, 0, 0, 0)
			}
		}
		if sp := a.assets.Sprites[name]; sp != nil && len(sp.Frames) > 0 {
			return sp.Frames[(a.title.tick/6+i)%len(sp.Frames)]
		}
		return nil
	}
	put := func(img *image.RGBA, x, y int) {
		if img != nil {
			draw.Draw(dst, img.Bounds().Add(image.Pt(x, y)), img, img.Bounds().Min, draw.Over)
		}
	}
	for x := 0; x < game.ScreenW; x += 32 { // the floor
		if a.skin != nil && len(a.skin.BlockTop) > 0 {
			put(a.skin.BlockTop[0], x, 448)
		} else {
			put(frame("sprite_block", 0), x, 448)
		}
	}
	put(frame("sprite_right", 0), 72, 416)
	put(frame("sprite_monsterl", 1), 520, 416)
	put(frame("sprite_flyingl", 2), 452, 360)
	put(frame("sprite_bonus", 3), 140, 400)
	put(frame("diamanto", 4), 576, 300)
	for i := range 5 {
		put(frame("sprite_stelo_poentoj", i), 208+i*40, 404)
	}
	o := &a.overlay
	o.AddRect(game.Rect{X: 0, Y: 0, W: game.ScreenW, H: game.ScreenH, Color: color.RGBA{0, 0, 0, 70}})
	if a.title.info != "" { // an info page instead of the title's texts (the overlay draws texts over every rect)
		a.drawInfo(o, a.title.info, frame)
		return
	}
	shadow := color.RGBA{0x3b, 0x22, 0x0c, 255}
	o.Add(game.Text{S: "Stelo", X: game.ScreenW/2 + 4, Y: 54, Size: 96, Weight: game.Black, Color: shadow, Align: game.Centre})
	o.Add(game.Text{S: "Stelo", X: game.ScreenW / 2, Y: 50, Size: 96, Weight: game.Black, Color: color.RGBA{0xff, 0xd0, 0x3a, 255}, Align: game.Centre})
	o.Add(game.Text{S: game.Tr("Platforma ludo · 2006, refarita 2026"), X: game.ScreenW / 2, Y: 186, Size: 15, Weight: game.Bold, Color: color.RGBA{0xff, 0xf4, 0xd8, 255}, Align: game.Centre})
	a.title.menu.Draw(o, game.MenuColors{Text: color.RGBA{0xff, 0xf4, 0xd8, 255}, Disabled: color.RGBA{0xb0, 0xa0, 0x80, 255},
		Selected: color.RGBA{0x2a, 0x1c, 0x12, 255}, Highlight: color.RGBA{0xff, 0xd0, 0x3a, 230}})
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
