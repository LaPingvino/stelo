// Command stelo is the Stelo remake (Ebitengine; desktop and WebAssembly).
package main

import (
	"image"
	"log"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/game"
	"github.com/LaPingvino/stelo/internal/gm"
	"github.com/LaPingvino/stelo/internal/store"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// The original's rooms in play order (the bonus level is reached in play).
var order = game.RoomOrder

type app struct {
	assets *game.Assets
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
	g.Highscores = store.LoadHighscores()
	g.SaveHighscores = store.SaveHighscores
	a.game = g
	ebiten.SetWindowTitle("Stelo")
}

func (a *app) Update() error {
	// Development aid until the title and cheats are in: Page Up/Down restart at
	// the previous/next room.
	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		a.load(a.room + 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		a.load(a.room - 1)
	}
	if a.game != nil {
		if err := a.game.Step(keys{}); err != nil {
			log.Print(err)
		}
	}
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	if a.game == nil {
		return
	}
	a.game.Draw(a.frame)
	screen.WritePixels(a.frame.Pix)
}

func (a *app) Layout(int, int) (int, int) { return game.ScreenW, game.ScreenH }

func main() {
	assets, err := game.LoadAssets(stelo.Files)
	if err != nil {
		log.Fatal(err)
	}
	a := &app{assets: assets, frame: image.NewRGBA(image.Rect(0, 0, game.ScreenW, game.ScreenH))}
	a.load(0)         // the title, then the rules screen
	ebiten.SetTPS(30) // every room runs at 30 steps per second
	ebiten.SetWindowSize(game.ScreenW, game.ScreenH)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(a); err != nil {
		log.Fatal(err)
	}
}

// keys is the keyboard as the game's Input (the player's keys come in with
// the rules, LOOP item 5).
type keys struct{}

func (keys) Held(k gm.Key) bool    { return ebiten.IsKeyPressed(ebitenKey(k)) }
func (keys) Pressed(k gm.Key) bool { return inpututil.IsKeyJustPressed(ebitenKey(k)) }

// Typed is what was typed since the last frame (keyboard_string; the cheat
// codes).
func (keys) Typed() []rune { return ebiten.AppendInputChars(nil) }

// AnyPressed is keyboard_wait's "any key".
func (keys) AnyPressed() bool { return len(inpututil.AppendJustPressedKeys(nil)) > 0 }

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
