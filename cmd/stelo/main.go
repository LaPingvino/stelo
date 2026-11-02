// Command stelo is the Stelo remake (Ebitengine; desktop and WebAssembly).
package main

import (
	"fmt"
	"image"
	"log"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/game"
	"github.com/LaPingvino/stelo/internal/gm"
	"github.com/LaPingvino/stelo/internal/pack"
	"github.com/LaPingvino/stelo/internal/sound"
	"github.com/LaPingvino/stelo/internal/store"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// The original's rooms in play order (the bonus level is reached in play).
var order = game.RoomOrder

type app struct {
	assets *game.Assets
	touch  touchState
	edit   editState
	// packs are the playable packs: the original first, then the ones in
	// the config folder; pack is the one chosen on the title.
	packs   []*pack.Pack
	pack    int
	sound   *sound.Player
	game    *game.Game
	room    int
	frame   *image.RGBA
	canvas  *ebiten.Image // the frame as an image, scaled onto the screen
	overlay game.Overlay  // texts drawn sharp over it (text.go)
	screen  *ebiten.Image
	skin    *game.Skin // the drawing skin (nil: the 2006 art); F3 / ?original
	title   titleScreen
	pause   pauseMenu
}

func (a *app) load(i int) {
	a.room = (i + len(order)) % len(order)
	g, err := game.NewGame(stelo.Files, a.assets, order[a.room])
	if err != nil {
		log.Printf("level %s: %v", order[a.room], err)
		return
	}
	if a.sound == nil {
		a.sound = sound.NewPlayer(stelo.Files)
		theSound = a.sound
		a.sound.SetMuted(store.Muted())
	}
	g.G.Sound = a.sound.Play
	g.Highscores = store.LoadHighscores()
	g.SaveHighscores = store.SaveHighscores
	g.OnLevel = func(n int) { store.Reached("Stelo", n) } // the original pack's progress
	a.game = g
	a.applySkin()
	ebiten.SetWindowTitle("Stelo")
}

func (a *app) Update() error {
	if files := ebiten.DroppedFiles(); files != nil { // a level file dropped on the window
		a.openDropped(files)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) { // the original art ⇄ the new skin
		a.toggleSkin()
		if a.title.on {
			a.buildTitleMenu(a.title.menu.Sel)
		}
		return nil
	}
	// The two corner buttons by mouse (no touch layout), and M for the sound.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !a.touch.seen {
		pt := image.Pt(ebiten.CursorPosition())
		switch {
		case pt.In(a.touch.lay.full):
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
			return nil
		case pt.In(a.touch.lay.snd):
			toggleSound()
			return nil
		case a.touch.showPause && pt.In(a.touch.lay.paus):
			if a.pause.on {
				a.pause.on = false
			} else {
				a.openPause()
			}
			return nil
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) && !a.edit.on {
		toggleSound()
	}
	// Full screen: F11, or Alt+Enter (the touch controls have a button).
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) || (inpututil.IsKeyJustPressed(ebiten.KeyEnter) && ebiten.IsKeyPressed(ebiten.KeyAlt)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
		return nil
	}
	if a.title.on {
		a.touch.update() // the controls appear (and full screen works) on the title too
		a.updateTitle()
		return nil
	}
	if !a.edit.on && a.game == nil { // left the editor with no game running
		a.openTitle()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) && a.game != nil { // the editor, on this room
		if a.edit.on {
			a.edit.on = false // back to the game as it was
		} else if a.game.World != nil {
			a.openEditor(a.game.World.Level.Name)
		}
		return nil
	}
	if a.edit.on {
		a.updateEditor()
		return nil
	}
	// The title's pack line (Pits' way): Left/Right choose while the title
	// waits; the next other key starts the chosen pack.
	if a.game != nil && a.game.Title && len(a.packs) > 1 {
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft):
			a.pack = (a.pack + len(a.packs) - 1) % len(a.packs)
			return nil
		case inpututil.IsKeyJustPressed(ebiten.KeyArrowRight):
			a.pack = (a.pack + 1) % len(a.packs)
			return nil
		}
		if a.pack > 0 && (len(inpututil.AppendJustPressedKeys(nil)) > 0 || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)) {
			a.playPack(a.packs[a.pack])
			return nil
		}
	}
	// Development aid until the title and cheats are in: Page Up/Down restart at
	// the previous/next room.
	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		a.load(a.room + 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		a.load(a.room - 1)
	}
	if a.game != nil {
		a.touch.update()
		a.touch.showPause = a.game.Level() > 0
		if a.pause.on {
			a.updatePause()
			return nil
		}
		if a.pauseWanted() {
			a.openPause()
			return nil
		}
		if err := a.game.Step(keys{&a.touch, pollPhoneKeys()}); err != nil {
			log.Print(err)
		}
	}
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	if a.title.on {
		a.overlay.Reset()
		a.drawTitle()
		a.canvas.WritePixels(a.frame.Pix)
		present(screen, a.canvas, a.frameArea(screen))
		drawOverlay(screen, &a.overlay)
		a.touch.draw(screen)
		return
	}
	if a.edit.on {
		a.edit.screen.Draw(a.frame)
		if a.game != nil {
			a.edit.screen.DrawPrompt(a.frame, a.edit.prompt, a.game.TextFont)
		}
		a.canvas.WritePixels(a.frame.Pix)
		present(screen, a.canvas, a.frameArea(screen))
		return
	}
	if a.game == nil {
		return
	}
	a.overlay.Reset()
	if a.game.World != nil {
		a.game.World.Overlay = &a.overlay // the HUD's text, drawn sharp after the frame (text.go)
	}
	a.game.Draw(a.frame)
	if a.game.Title && len(a.packs) > 1 {
		p := a.packs[a.pack]
		label := p.Name
		if n := store.Progress(p.Name); n > 0 { // the furthest level reached in it
			label = fmt.Sprintf("%s (%d/%d)", p.Name, min(n, len(p.Levels)), len(p.Levels))
		}
		a.game.DrawPackChoice(a.frame, label)
	}
	a.canvas.WritePixels(a.frame.Pix)
	present(screen, a.canvas, a.frameArea(screen))
	if a.pause.on {
		a.drawPause()
	}
	drawOverlay(screen, &a.overlay)
	a.touch.draw(screen)
}

// Layout is the screen's real resolution; the frame is scaled onto it (scale.go).
func (a *app) Layout(outW, outH int) (int, int) { return screenSize(outW, outH) }

func main() {
	assets, err := game.LoadAssets(stelo.Files)
	if err != nil {
		log.Fatal(err)
	}
	a := &app{assets: assets, frame: image.NewRGBA(image.Rect(0, 0, game.ScreenW, game.ScreenH)), canvas: ebiten.NewImage(game.ScreenW, game.ScreenH)}
	loadNoto(stelo.Files)
	a.touch.seen = showPadsFromURL()
	if o, err := pack.Original(stelo.Files); err == nil {
		a.packs = append([]*pack.Pack{o}, store.ListPacks()...)
	}
	if langFromURL() == "en" { // ?lang=en / STELO_LANG=en: the English translation
		game.Lang = "en"
	}
	if !originalArtFromURL() { // the new art by default (Joop 2026-11-02); ?original / STELO_ART=original: the 2006 art
		a.toggleSkin()
	}
	a.openTitle()      // the new title (docs/UI.md); the 2006 title picture and rules room are no longer shown
	if editFromURL() { // ?edit: straight into the editor on level 1
		a.title.on = false
		a.openEditor("nivelo1")
	}
	ebiten.SetTPS(30) // every room runs at 30 steps per second
	ebiten.SetWindowSize(game.ScreenW, game.ScreenH)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(a); err != nil {
		log.Fatal(err)
	}
}

// keys is the keyboard as the game's Input (the player's keys come in with
// the rules, LOOP item 5).
type keys struct {
	t     *touchState
	phone phoneKeys // this frame's characters and Enter from the phone keyboard
}

func (k keys) Held(key gm.Key) bool {
	return ebiten.IsKeyPressed(ebitenKey(key)) || (k.t != nil && k.t.held[key])
}
func (k keys) Pressed(key gm.Key) bool {
	return inpututil.IsKeyJustPressed(ebitenKey(key)) || (key == gm.KeyEnter && k.phone.enter)
}

// Typed is what was typed since the last frame (keyboard_string; the cheat
// codes).
func (k keys) Typed() []rune {
	t := ebiten.AppendInputChars(nil)
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		t = append(t, '\b')
	}
	return append(t, k.phone.chars...)
}

// AnyPressed is keyboard_wait's "any key".
func (k keys) AnyPressed() bool {
	return len(inpututil.AppendJustPressedKeys(nil)) > 0 || (k.t != nil && k.t.tapped) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || len(k.phone.chars) > 0 || k.phone.enter
}

func ebitenKey(k gm.Key) ebiten.Key {
	switch k {
	case gm.KeyLeft:
		return ebiten.KeyArrowLeft
	case gm.KeyRight:
		return ebiten.KeyArrowRight
	case gm.KeyUp:
		return ebiten.KeyArrowUp
	case gm.KeyDown:
		return ebiten.KeyArrowDown
	case gm.KeySpace:
		return ebiten.KeySpace
	case gm.KeyEnter:
		return ebiten.KeyEnter
	case gm.KeyPgDn:
		return ebiten.KeyPageDown
	}
	return ebiten.KeyF24
}

// toggleSkin switches between the original art and the new skin (F3).
func (a *app) toggleSkin() {
	if a.skin == nil {
		s, err := game.LoadNovaSkin(stelo.Files)
		if err != nil {
			log.Printf("skin: %v", err)
			return
		}
		a.skin = s
	} else {
		a.skin = nil
	}
	a.applySkin()
}

// applySkin gives the running game the app's skin (nil: the 2006 art).
func (a *app) applySkin() {
	if a.game == nil {
		return
	}
	a.game.Skin = a.skin
	if a.game.World != nil {
		a.game.World.Skin = a.skin
	}
}

// frameArea is where the frame goes on the screen: all of it, or — once touch
// has been used — what the touch controls leave (touch.go); it also stores the
// layout the touches are read against.
func (a *app) frameArea(screen *ebiten.Image) image.Rectangle {
	b := screen.Bounds()
	if !a.touch.seen { // no touch: the whole screen; the two buttons small in the top-left corner
		d := ebiten.Monitor().DeviceScaleFactor()
		k := int(28 * max(d, 1))
		m := k / 4
		a.touch.lay = touchLayout{area: b, full: image.Rect(m, m, m+k, m+k*2/3), snd: image.Rect(2*m+k, m, 2*m+2*k, m+k*2/3),
			paus: image.Rect(3*m+2*k, m, 3*m+3*k, m+k*2/3)}
		return b
	}
	a.touch.lay = layoutTouch(b.Dx(), b.Dy())
	return a.touch.lay.area
}

// theSound is the program's sound player (made once).
var theSound *sound.Player

// toggleSound is the sound button and M: all sound off or on, remembered.
func toggleSound() {
	if theSound == nil {
		return
	}
	theSound.SetMuted(!theSound.Muted())
	store.SetMuted(theSound.Muted())
}

func soundMuted() bool { return theSound != nil && theSound.Muted() }
