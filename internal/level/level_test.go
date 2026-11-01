package level

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Every original room survives Format → Parse unchanged (instances compared
// as a set: the file doesn't keep creation order), and levels/ is up to date.
func TestOriginalRoundTrip(t *testing.T) {
	levels, err := LoadOriginal("../../original/project.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(levels) != 9 {
		t.Fatalf("%d rooms, want 9", len(levels))
	}
	for _, want := range levels {
		text := Format(want)
		got, err := Parse(text)
		if err != nil {
			t.Fatalf("%s: %v", want.Name, err)
		}
		w, g := *want, *got
		w.Instances = append([]Instance(nil), want.Instances...)
		SortInstances(w.Instances)
		SortInstances(g.Instances)
		if !reflect.DeepEqual(w, g) {
			t.Errorf("%s: round trip differs", want.Name)
		}
		file, err := os.ReadFile(filepath.Join("../../levels", want.Name+".txt"))
		if err != nil || string(file) != text {
			t.Errorf("levels/%s.txt is missing or out of date: go run ./cmd/convert", want.Name)
		}
	}
}
