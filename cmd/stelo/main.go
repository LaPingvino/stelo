// Command stelo is the Stelo remake (Ebitengine; desktop and WebAssembly).
package main

import (
	"image"
	"log"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/game"
	"github.com/LaPingvino/stelo/internal/gm"
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
	sound  *sound.Player
	game   *game.Game
	room   int
	frame  *image.RGBA
	screen *ebiten.Image
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
	}
	g.G.Sound = a.sound.Play
	g.Highscores = store.LoadHighscores()
	g.SaveHighscores = store.SaveHighscores
	a.game = g
	ebiten.SetWindowTitle("Stelo")
}

func (a *app) Update() error {
	if files := ebiten.DroppedFiles(); files != nil { // a level file dropped on the window
		a.openDropped(files)
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
		if err := a.game.Step(keys{&a.touch, pollPhoneKeys()}); err != nil {
			log.Print(err)
		}
	}
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	if a.edit.on {
		a.edit.screen.Draw(a.frame)
		screen.WritePixels(a.frame.Pix)
		return
	}
	if a.game == nil {
		return
	}
	a.game.Draw(a.frame)
	screen.WritePixels(a.frame.Pix)
	a.touch.draw(screen)
}

func (a *app) Layout(int, int) (int, int) { return game.ScreenW, game.ScreenH }

func main() {
	assets, err := game.LoadAssets(stelo.Files)
	if err != nil {
		log.Fatal(err)
	}
	a := &app{assets: assets, frame: image.NewRGBA(image.Rect(0, 0, game.ScreenW, game.ScreenH))}
	a.load(0) // the title, then the rules screen
	a.touch.seen = showPadsFromURL()
	if editFromURL() { // ?edit: straight into the editor on level 1
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
