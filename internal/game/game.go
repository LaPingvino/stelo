package game

import (
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
	FS     fs.FS
	Assets *Assets
	G      *gm.Globals
	World  *World
	saved  *World // the persistent level left for the bonus level
}

// NewGame starts at a room (normally RoomOrder[0]).
func NewGame(fsys fs.FS, a *Assets, room string) (*Game, error) {
	g := &Game{FS: fsys, Assets: a, G: &gm.Globals{Lives: 3, Vars: map[string]float64{}}}
	return g, g.enter(room)
}

func (g *Game) enter(room string) error {
	w, err := LoadLevel(g.FS, g.Assets, room)
	if err != nil {
		return err
	}
	w.Room.G = g.G
	g.World = w
	return nil
}

// Step runs one step and then any room change it asked for.
func (g *Game) Step(in gm.Input) error {
	g.World.Step(in)
	if g.G.NoMoreLives { // life_controller: "no more lives" → finlisto
		g.G.NoMoreLives = false
		g.World.Room.Goto = "finlisto"
	}
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
