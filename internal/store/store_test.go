//go:build !js

package store

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/pack"

	"github.com/LaPingvino/stelo/internal/game"
)

func TestHighscoresRoundTrip(t *testing.T) {
	dir := t.TempDir()
	Dir = func() string { return dir }
	if h := LoadHighscores(); h != game.ClearedHighscores() {
		t.Fatalf("nothing saved yet: %+v", h)
	}
	h := game.ClearedHighscores()
	h.Add("Joop", 1230)
	SaveHighscores(h)
	if got := LoadHighscores(); got != h {
		t.Errorf("loaded %+v, want %+v", got, h)
	}
}

func TestLevelRoundTrip(t *testing.T) {
	dir := t.TempDir()
	Dir = func() string { return dir }
	if _, ok := LoadLevel("nivelo1"); ok {
		t.Fatal("nothing saved yet")
	}
	SaveLevel("nivelo1", "stelo-level 1\nname nivelo1\n")
	if got, ok := LoadLevel("nivelo1"); !ok || got != "stelo-level 1\nname nivelo1\n" {
		t.Errorf("loaded %q %v", got, ok)
	}
}

func TestSavePackUnzips(t *testing.T) {
	dir := t.TempDir()
	Dir = func() string { return dir }
	var buf bytes.Buffer
	z := zip.NewWriter(&buf)
	w, _ := z.Create("Mia pako/pack.txt")
	w.Write([]byte("name=Mia pako\n"))
	z.Close()
	SavePack("Mia pako", buf.Bytes())
	data, err := os.ReadFile(filepath.Join(dir, "packs", "Mia pako", "pack.txt"))
	if err != nil || string(data) != "name=Mia pako\n" {
		t.Errorf("%q %v", data, err)
	}
}

// A pack saved with SavePack is listed.
func TestListPacks(t *testing.T) {
	dir := t.TempDir()
	Dir = func() string { return dir }
	o, err := pack.Original(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	o.Name = "Testo"
	data, _ := o.Zip()
	SavePack(o.Name, data)
	ps := ListPacks()
	if len(ps) != 1 || ps[0].Name != "Testo" || len(ps[0].Levels) != 5 {
		t.Fatalf("%d packs", len(ps))
	}
}
