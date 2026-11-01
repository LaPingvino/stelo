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

// The exit goes to the next level after the 1 s pause; the diamond goes to
// the bonus level and a flyer there brings the player back to the saved level
// as it was left.
func TestRoomChanges(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(stelo.Files, a, "nivelo1")
	if err != nil {
		t.Fatal(err)
	}
	p := g.World.Room.All("character")[0]
	// The diamond.
	d := g.World.Room.All("iru_al_bonuso")[0]
	p.X, p.Y = d.X, d.Y
	g.Step(held{})
	if g.World.Level.Name != "bonusnivelo" {
		t.Fatalf("diamond: in %s", g.World.Level.Name)
	}
	// A flyer in the bonus level: back to the saved level, no life lost.
	lives := g.G.Lives
	bp := g.World.Room.All("character")[0]
	f := g.World.Room.All("flyer")[0]
	bp.X, bp.Y = f.X, f.Y
	for range 40 {
		g.Step(held{})
	}
	if g.World.Level.Name != "nivelo1" || g.G.Lives != lives {
		t.Fatalf("flyer in the bonus level: in %s, lives %d (was %d)", g.World.Level.Name, g.G.Lives, lives)
	}
	if len(g.World.Room.All("iru_al_bonuso")) != 0 || g.World.Room.All("character")[0] != p {
		t.Error("the saved level is the same room, diamond taken")
	}
	// The exit: next room after the pause.
	e := g.World.Room.All("levelexit")[0]
	p.X, p.Y = e.X, e.Y
	for range 40 {
		g.Step(held{})
	}
	if g.World.Level.Name != "nivelo2" {
		t.Errorf("exit: in %s", g.World.Level.Name)
	}
}

// Level 5's spawners: with no monsters left, each starts a wave — monsters at
// steps 30, 60 and 90 at the spawner — and may start again after the third.
func TestLevel5Waves(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(stelo.Files, a, "nivelo5")
	if err != nil {
		t.Fatal(err)
	}
	r := g.World.Room
	spawners := len(r.All("finmonsterfarilo"))
	for _, m := range r.All("monster") {
		m.Destroyed = true
	}
	g.World.Room.All("character")[0].Vars["netusxu"] = 1 // keep the player out of it
	count := func() int { return len(r.All("monster")) }
	for range 31 {
		g.Step(held{})
	}
	if count() != spawners {
		t.Fatalf("after 31 steps: %d monsters, want one per spawner (%d)", count(), spawners)
	}
	for range 60 {
		g.Step(held{})
	}
	if count() != 3*spawners {
		t.Errorf("after 91 steps: %d monsters, want %d", count(), 3*spawners)
	}
}

// The bonus level's b_monsters fall (gravity) and land on blocks.
func TestBonusMonstersFall(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(stelo.Files, a, "bonusnivelo")
	if err != nil {
		t.Fatal(err)
	}
	g.World.Room.All("character")[0].Vars["netusxu"] = 1
	for range 31 {
		g.Step(held{})
	}
	bs := g.World.Room.All("b_monster")
	if len(bs) == 0 {
		t.Fatal("no b_monster after 31 steps")
	}
	y0 := bs[0].Y
	for range 60 {
		g.Step(held{})
	}
	if bs[0].Y <= y0 {
		t.Errorf("b_monster did not fall: y %v → %v", y0, bs[0].Y)
	}
}

// Every room loads and runs 10 seconds with no keys held.
func TestEveryRoomRuns(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	for _, room := range RoomOrder {
		g, err := NewGame(stelo.Files, a, room)
		if err != nil {
			t.Fatalf("%s: %v", room, err)
		}
		for range 300 {
			if err := g.Step(held{}); err != nil {
				t.Fatalf("%s: %v", room, err)
			}
		}
	}
}

// danklisto: the credits scroll as the view follows the b_monster zig-zagging
// down the 960 px room; the life controller (HUD) is gone.
func TestCreditsScroll(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(stelo.Files, a, "danklisto")
	if err != nil {
		t.Fatal(err)
	}
	if g.G.Vars["life_controller_gone"] != 1 {
		t.Error("danklisto's creation code destroys the life controller")
	}
	y0 := g.World.ViewY
	// The b_monster zig-zags down floor by floor (~500 px a floor at 4 px a
	// step); the view starts moving once it is within 128 px of the bottom.
	for range 1200 {
		g.Step(held{})
	}
	if g.World.ViewY <= y0 {
		t.Errorf("view did not scroll down: %d → %d", y0, g.World.ViewY)
	}
}

type anyKey struct {
	held
	any bool
}

func (a anyKey) AnyPressed() bool { return a.any }

// Game start: the title picture until a key, the rules screen until a key,
// then level 1 (klarigilo).
func TestTitleAndRules(t *testing.T) {
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(stelo.Files, a, RoomOrder[0])
	if err != nil {
		t.Fatal(err)
	}
	for range 30 {
		g.Step(anyKey{})
	}
	if !g.Title || g.World.Level.Name != "klarigo" {
		t.Fatalf("waits on the title: title %v in %s", g.Title, g.World.Level.Name)
	}
	g.Step(anyKey{any: true})
	if g.Title || !g.Waiting() || g.World.Level.Name != "klarigo" {
		t.Fatalf("after a key: title %v waiting %v in %s", g.Title, g.Waiting(), g.World.Level.Name)
	}
	g.Step(anyKey{any: true})
	if g.World.Level.Name != "nivelo1" || g.Waiting() {
		t.Fatalf("after a second key: in %s, waiting %v", g.World.Level.Name, g.Waiting())
	}
}
