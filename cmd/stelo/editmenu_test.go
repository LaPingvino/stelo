package main

import (
	"testing"

	"github.com/LaPingvino/stelo/internal/level"
)

// The editor menu's new level parses: 20×15, walls round the edge, the
// player and the exit inside.
func TestNewLevelText(t *testing.T) {
	l, err := level.Parse(newLevelText("nova1"))
	if err != nil {
		t.Fatal(err)
	}
	if l.Name != "nova1" {
		t.Errorf("name %q", l.Name)
	}
	var player, exit bool
	for _, in := range l.Instances {
		player = player || in.Object == "character"
		exit = exit || in.Object == "levelexit"
	}
	if !player || !exit {
		t.Errorf("player %v, exit %v", player, exit)
	}
}
