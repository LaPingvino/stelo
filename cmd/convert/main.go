// Command convert writes the original rooms (original/project.json) as level
// files in levels/ (docs/LEVEL-FORMAT.md).
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/LaPingvino/stelo/internal/level"
)

func main() {
	levels, err := level.LoadOriginal("original/project.json")
	if err != nil {
		log.Fatal(err)
	}
	for _, l := range levels {
		p := filepath.Join("levels", l.Name+".txt")
		if err := os.WriteFile(p, []byte(level.Format(l)), 0o644); err != nil {
			log.Fatal(err)
		}
		log.Printf("%s: %d instances", p, len(l.Instances))
	}
}
