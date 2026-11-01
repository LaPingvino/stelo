// Command stelo is the Stelo remake (Ebitengine; desktop and WebAssembly).
package main

import (
	"image"
	"log"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/game"
	"github.com/LaPingvino/stelo/internal/gm"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// The original's rooms in play order (the bonus level is reached in play).
var order = []string{"klarigo", "nivelo1", "nivelo2", "nivelo3", "nivelo4", "nivelo5", "danklisto", "finlisto", "bonusnivelo"}

type app struct {
	assets *game.Assets
	world  *game.World
	room   int
	frame  *image.RGBA
	screen *ebiten.Image
}

func (a *app) load(i int) {
	a.room = (i + len(order)) % len(order)
	w, err := game.LoadLevel(stelo.Files, a.assets, order[a.room])
	if err != nil {
		log.Printf("level %s: %v", order[a.room], err)
		return
	}
	a.world = w
	ebiten.SetWindowTitle("Stelo — " + order[a.room])
}

func (a *app) Update() error {
	// Until the rules are in (LOOP item 4+): Page Up/Down step through the rooms.
	if inpututil.IsKeyJustPressed(ebiten.KeyPageDown) {
		a.load(a.room + 1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyPageUp) {
		a.load(a.room - 1)
	}
	if a.world != nil {
		a.world.Step(keys{})
	}
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	if a.world == nil {
		return
	}
	a.world.Draw(a.frame)
	screen.WritePixels(a.frame.Pix)
}

func (a *app) Layout(int, int) (int, int) { return game.ScreenW, game.ScreenH }

func main() {
	assets, err := game.LoadAssets(stelo.Files)
	if err != nil {
		log.Fatal(err)
	}
	a := &app{assets: assets, frame: image.NewRGBA(image.Rect(0, 0, game.ScreenW, game.ScreenH))}
	a.load(1)
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
