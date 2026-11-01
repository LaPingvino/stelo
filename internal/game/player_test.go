package game

import (
	"testing"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/gm"
)

type held map[gm.Key]bool

func (h held) Held(k gm.Key) bool    { return h[k] }
func (h held) Pressed(k gm.Key) bool { return false }

func level1(t *testing.T) (*World, *gm.Instance) {
	t.Helper()
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	w, err := LoadLevel(stelo.Files, a, "nivelo1")
	if err != nil {
		t.Fatal(err)
	}
	p := w.Room.All("character")
	if len(p) != 1 {
		t.Fatalf("%d players", len(p))
	}
	return w, p[0]
}

// Level 1: the player starts on the floor at (32, 416) (cell 1, row 13). A
// ladder stands at cell 2 of that row (box x 70–89).
func TestPlayerStandsJumpsWalksAndClimbs(t *testing.T) {
	w, p := level1(t)
	for range 10 {
		w.Step(held{})
	}
	if p.X != 32 || p.Y != 416 {
		t.Fatalf("standing still moved the player to %v,%v", p.X, p.Y)
	}
	// Jump at the start (box x 39–55: clear of the ladder).
	w.Step(held{gm.KeyUp: true})
	if p.VSpeed >= 0 || p.Y >= 416 {
		t.Fatalf("jump: y %v vspeed %v", p.Y, p.VSpeed)
	}
	top := p.Y
	for range 60 {
		w.Step(held{})
		top = min(top, p.Y)
	}
	if p.Y != 416 || p.VSpeed != 0 {
		t.Errorf("landed at y %v vspeed %v, want 416 and 0", p.Y, p.VSpeed)
	}
	// vspeed −10 with gravity 0.5 from the next step on: 10 + 9.5 + … ≈ 105 px.
	if rise := 416 - top; rise < 95 || rise > 110 {
		t.Errorf("jump height %v px", rise)
	}
	// Walk right: 4 px a step, facing right…
	w.Step(held{gm.KeyRight: true})
	if p.X != 36 || p.Sprite != "sprite_right" {
		t.Errorf("one step right: x %v %s", p.X, p.Sprite)
	}
	// …until the box reaches the ladder: the step event then shows the
	// climbing sprite after the key event set sprite_right (as the original).
	for range 4 {
		w.Step(held{gm.KeyRight: true})
	}
	if p.X != 52 || p.Y != 416 || p.Sprite != "sprite_climbing" {
		t.Errorf("after 5 steps right: %v,%v %s; want 52,416 sprite_climbing", p.X, p.Y, p.Sprite)
	}
	// Now overlapping the ladder: Up climbs 3 px a step, no gravity.
	for range 4 {
		w.Step(held{gm.KeyUp: true})
	}
	if p.Y != 404 || p.Sprite != "sprite_climbing" {
		t.Errorf("climbing: y %v sprite %s, want 404 sprite_climbing", p.Y, p.Sprite)
	}
	w.Step(held{})
	if p.Y != 404 {
		t.Errorf("hangs on the ladder when no key is held: y %v", p.Y)
	}
}

// A star coin is worth 10; stomping a monster from above is worth 50 and
// leaves a flattened monster for 10 steps; walking into one costs a life.
func TestCoinsAndMonsters(t *testing.T) {
	w, p := level1(t)
	g := w.Room.G
	a := w.Assets
	coin := w.Room.Create(a.Object("stelo_poentoj"), p.X, p.Y)
	w.Step(held{})
	if g.Score != 10 || !coin.Destroyed {
		t.Fatalf("coin: score %d", g.Score)
	}
	// A monster just below the falling player: stomp.
	m := w.Room.Create(a.Object("monster"), p.X, p.Y)
	m.HSpeed = 0
	p.Y -= 20
	p.VSpeed = 4
	m.Y = p.Y + 10 // player above m.Y+8 while falling
	w.Step(held{})
	if g.Score != 60 || m.Obj.Name != "monster_dead" {
		t.Fatalf("stomp: score %d, monster is %s", g.Score, m.Obj.Name)
	}
	for range 10 {
		w.Step(held{})
	}
	if !m.Destroyed {
		t.Error("flattened monster gone after 10 steps")
	}
	// A monster walked into from the side: a life less, back to the start.
	for range 40 {
		w.Step(held{})
	}
	lives := g.Lives
	m2 := w.Room.Create(a.Object("monster"), p.X, p.Y)
	m2.HSpeed = 0
	w.Step(held{})
	if g.Lives != lives-1 || w.Room.Sleep == 0 {
		t.Errorf("hit: lives %d (was %d), sleep %d", g.Lives, lives, w.Room.Sleep)
	}
}
