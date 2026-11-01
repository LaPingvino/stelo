package pack

import (
	"archive/zip"
	"bytes"
	"testing"
	"testing/fstest"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/level"
)

func TestLoad(t *testing.T) {
	l1, _ := stelo.Files.ReadFile("levels/nivelo1.txt")
	l3, _ := stelo.Files.ReadFile("levels/nivelo3.txt")
	fsys := fstest.MapFS{
		"packs/mine/pack.txt": {Data: []byte("name=Mia unua pako\nauthor = Joop\n")},
		"packs/mine/02.txt":   {Data: l1},
		"packs/mine/01.txt":   {Data: l3},
		"packs/mine/notes.md": {Data: []byte("x")},
		"packs/bare/a.txt":    {Data: l1},
		"packs/broken/a.txt":  {Data: []byte("nope")},
	}
	p, err := Load(fsys, "packs/mine")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Mia unua pako" || p.Author != "Joop" || len(p.Levels) != 2 || p.Levels[0].Name != "nivelo3" {
		t.Errorf("%q by %q: %d levels, first %s", p.Name, p.Author, len(p.Levels), p.Levels[0].Name)
	}
	if p, err := Load(fsys, "packs/bare"); err != nil || p.Name != "bare" {
		t.Errorf("a pack without pack.txt: %v %v", p, err)
	}
	if _, err := Load(fsys, "packs/broken"); err == nil {
		t.Error("a broken level wasn't reported")
	}
	o, err := Original(stelo.Files)
	if err != nil || len(o.Levels) != 5 || o.Levels[4].Name != "nivelo5" {
		t.Errorf("original pack: %v", err)
	}
}

// A pack written as a zip reads back as the same pack (unzipped: a folder).
func TestZipRoundTrip(t *testing.T) {
	o, err := Original(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	o.Name, o.Author = "Mia/pako", "Joop"
	data, err := o.Zip()
	if err != nil {
		t.Fatal(err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	p, err := Load(zr, "Mia_pako")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Mia/pako" || p.Author != "Joop" || len(p.Levels) != 5 {
		t.Fatalf("%q by %q, %d levels", p.Name, p.Author, len(p.Levels))
	}
	for i := range p.Levels {
		if level.Format(p.Levels[i]) != level.Format(o.Levels[i]) {
			t.Errorf("level %d changed", i+1)
		}
	}
}
