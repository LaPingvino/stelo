package game

import (
	"strings"

	"github.com/LaPingvino/stelo/internal/gm"
)

// cheat is one of the character's cheat codes (the code action in its step
// event): when keyboard_string contains the code (string_count > 0, case
// sensitive), keyboard_string is cleared and the effect runs.
type cheat struct {
	code   string
	effect func(r *gm.Room, self *gm.Instance, a *Assets)
}

var cheats = []cheat{
	{"sekvosio", func(r *gm.Room, self *gm.Instance, a *Assets) { r.Goto = "@next" }}, // room_goto_next()
	{"vivuplilonge", func(r *gm.Room, self *gm.Instance, a *Assets) { r.G.Lives++ }},
	{"netusxumin", func(r *gm.Room, self *gm.Instance, a *Assets) { self.Vars["netusxu"] = 1 }},
	{"mibedauxras", func(r *gm.Room, self *gm.Instance, a *Assets) {
		self.Vars["netusxu"] = 0
		r.G.Play("@music:muziko1") // sound_stop_all(); sound_loop(muziko1)
	}},
	{"esperanto", func(r *gm.Room, self *gm.Instance, a *Assets) { r.G.Play("@music:LaEspero") }},
	{"fekalfajro", func(r *gm.Room, self *gm.Instance, a *Assets) {
		for _, d := range r.All("death") { // with (death) instance_change(block, false)
			r.Change(d, a.Object("block"), false)
		}
	}},
	{"miamasGrandanFraton", func(r *gm.Room, self *gm.Instance, a *Assets) { r.G.Score += 500 }},
	{"pasispasinteco", func(r *gm.Room, self *gm.Instance, a *Assets) {
		r.G.Vars["highscore_clear"] = 1 // highscore_clear() (the table: item 8, highscores)
		r.G.Score = -1
	}},
	{"pliajpoentoj", func(r *gm.Room, self *gm.Instance, a *Assets) {
		r.G.Play("@music:bonussono")
		r.Goto = "@bonus" // room_persistent=true; global.ludnivelo=room; room_goto(bonusnivelo)
	}},
}

// runCheats is the cheat part of the character's step event, in the original's
// order (each check sees the string the previous one may have cleared).
func runCheats(r *gm.Room, self *gm.Instance, a *Assets) {
	for _, c := range cheats {
		if strings.Contains(r.G.KeyboardString, c.code) {
			r.G.KeyboardString = ""
			c.effect(r, self, a)
		}
	}
}
