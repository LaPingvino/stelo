// Command render draws a level's first frame to a PNG (for checks).
//
//	go run ./cmd/render -level nivelo1 -out nivelo1.png
package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/game"
	"github.com/LaPingvino/stelo/internal/gm"
)

func main() {
	name := flag.String("level", "nivelo1", "level name")
	out := flag.String("out", "frame.png", "output PNG")
	steps := flag.Int("steps", 0, "steps to run first (no keys pressed)")
	flag.Parse()
	a, err := game.LoadAssets(stelo.Files)
	if err != nil {
		log.Fatal(err)
	}
	w, err := game.LoadLevel(stelo.Files, a, *name)
	if err != nil {
		log.Fatal(err)
	}
	for range *steps {
		w.Step(noKeys{})
	}
	img := image.NewRGBA(image.Rect(0, 0, game.ScreenW, game.ScreenH))
	w.Draw(img)
	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

type noKeys struct{}

func (noKeys) Held(gm.Key) bool    { return false }
func (noKeys) Pressed(gm.Key) bool { return false }
