package game

import (
	"golang.org/x/image/font"
	"io/fs"

	"github.com/LaPingvino/stelo/internal/gm"
)

// RoomOrder is the original's room order (project.json room_order); the
// bonus level is reached only through the diamond (or a cheat).
var RoomOrder = []string{"klarigo", "nivelo1", "nivelo2", "nivelo3", "nivelo4", "nivelo5", "danklisto", "finlisto", "bonusnivelo"}

// Game runs the rooms one after another with the state that lives across
// them: score, lives, and the level saved while the bonus level is played
// (the original makes that room persistent and returns to global.ludnivelo).
type Game struct {
	FS       fs.FS
	Assets   *Assets
	G        *gm.Globals
	World    *World
	saved    *World // the persistent level left for the bonus level
	Font     font.Face
	TextFont font.Face
	// Title: the title picture is up (game start); waitRoom: klarigo waits for
	// a key (klarigilo's keyboard_wait).
	Title    bool
	waitRoom bool
}

// NewGame starts at a room (normally RoomOrder[0]).
func NewGame(fsys fs.FS, a *Assets, room string) (*Game, error) {
	g := &Game{FS: fsys, Assets: a, G: &gm.Globals{Lives: 3, Vars: map[string]float64{}}}
	if f, err := HUDFont(fsys); err == nil {
		g.Font = f
	}
	if f, err := TextFont(fsys); err == nil {
		g.TextFont = f
	}
	// klarigilo's game start event (the first room holds the klarigilo):
	// the komenco sound and the title picture until a key.
	if room == RoomOrder[0] {
		g.Title = true
		g.G.Play("@sound:komenco")
	}
	return g, g.enter(room)
}

func (g *Game) enter(room string) error {
	w, err := LoadLevel(g.FS, g.Assets, room)
	if err != nil {
		return err
	}
	w.Room.G = g.G
	w.Font, w.TextFont = g.Font, g.TextFont
	g.World = w
	// The persistent life controller exists from the room that places it
	// (nivelo1) until it is destroyed.
	if len(w.Room.All("life_controller")) > 0 {
		g.G.Vars["life_controller_alive"] = 1
	}
	if code := roomCode[room]; code != nil {
		code(g)
	}
	if len(w.Room.All("klarigilo")) > 0 {
		g.waitRoom = true
	}
	return nil
}

// roomCode is the rooms' creation code, ported (levels/*.txt "code" lines).
var roomCode = map[string]func(g *Game){
	// danklisto: with (life_controller) {instance_destroy();}
	// sound_stop_all(); sound_loop(LaEspero);
	"danklisto": func(g *Game) {
		g.G.Vars["life_controller_gone"] = 1
		g.G.Play("@music:LaEspero")
	},
}

// AnyKey is the part of the input keyboard_wait needs: a key went down this
// step (any key).
type AnyKey interface{ AnyPressed() bool }

// Waiting reports a keyboard_wait in progress: the title picture (Title) or
// the rules screen. The game draws but doesn't step while waiting.
func (g *Game) Waiting() bool { return g.Title || g.waitRoom }

// Step runs one step and then any room change it asked for.
func (g *Game) Step(in gm.Input) error {
	// klarigilo (original/objects/klarigilo.gml): at game start the title
	// picture stelo_komenco and keyboard_wait(); in klarigo's room start
	// screen_redraw(); keyboard_wait(); room_goto_next().
	if g.Waiting() {
		if a, ok := in.(AnyKey); ok && a.AnyPressed() {
			if g.Title {
				g.Title = false
				g.G.Play("@music:muziko1") // end komenco, loop muziko1
				return nil
			}
			g.waitRoom = false
			g.World.Room.Goto = "@next"
			return g.change()
		}
		return nil
	}
	g.World.Step(in)
	if g.G.NoMoreLives { // life_controller: "no more lives" → finlisto
		g.G.NoMoreLives = false
		g.World.Room.Goto = "finlisto"
	}
	return g.change()
}

// change performs the room change the room asked for, if any.
func (g *Game) change() error {
	to := g.World.Room.Goto
	if to == "" || g.World.Room.Sleep > 0 {
		return nil // a room change waits until the pause is over
	}
	g.World.Room.Goto = ""
	switch to {
	case "@next":
		for i, n := range RoomOrder {
			if n == g.World.Level.Name && i+1 < len(RoomOrder) {
				return g.enter(RoomOrder[i+1])
			}
		}
		return nil
	case "@bonus": // the room becomes persistent: keep it as it is
		g.saved = g.World
		return g.enter("bonusnivelo")
	case "@saved":
		if g.saved != nil {
			g.World, g.saved = g.saved, nil
			return nil
		}
		return nil
	}
	return g.enter(to)
}

// characterRooms adds the character's room-changing collisions.
func characterRooms(a *Assets, o *gm.Object) {
	{
		c := o.Events.Collision
		// Collision with levelexit: sound, 1 s pause, the next room (if any).
		c["levelexit"] = func(r *gm.Room, self, other *gm.Instance) {
			if invulnerable(self) {
				return
			}
			r.G.Play("finishlevel")
			r.SleepMS(1000, roomSpeed)
			r.Goto = "@next"
		}
		// Collision with iru_al_bonuso (the diamond): the room is kept as it
		// is, the diamond disappears, the bonus music, the bonus level.
		c["iru_al_bonuso"] = func(r *gm.Room, self, other *gm.Instance) {
			other.Destroyed = true
			r.G.Play("@music:bonussono")
			r.Goto = "@bonus"
		}
	}
}
