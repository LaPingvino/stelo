package game

import "github.com/LaPingvino/stelo/internal/gm"

// The bonus level's monsters and the spawners of the bonus level and level 5
// (original/objects/b_monster.gml, monsterfarilo.gml, finmonsterfarilo.gml,
// timelines/faruMonsterojn.gml, faruFinmonsterojn.gml).

func init() {
	objectEvents["b_monster"] = bMonsterEvents
	objectEvents["monsterfarilo"] = func(a *Assets, o *gm.Object) { spawnerEvents(a, o, "b_monster") }
	objectEvents["finmonsterfarilo"] = func(a *Assets, o *gm.Object) { spawnerEvents(a, o, "monster") }
}

// bMonsterEvents: as monster (right at 4, sprite by direction, back + reverse
// at a marker), plus gravity like the player's; on a block: move to contact in
// the direction of motion (max 12), vspeed 0, reverse.
func bMonsterEvents(a *Assets, o *gm.Object) {
	monsterEvents(a, o)
	o.Events.Step = func(r *gm.Room, self *gm.Instance) {
		if r.PlaceFree(self, self.X, self.Y+1) {
			self.Gravity, self.GravityDirection = 0.5, 270
		} else {
			self.Gravity = 0
		}
	}
	o.Events.Collision["block"] = func(r *gm.Room, self, other *gm.Instance) {
		r.MoveContactSolid(self, self.Direction(), 12)
		self.VSpeed = 0
		self.HSpeed = -self.HSpeed
	}
}

// spawnerEvents: each step, when no instance of what it makes is left and its
// timeline isn't running, start the timeline: three of them at steps 30, 60
// and 90, at the spawner's position; after the third it may start again.
func spawnerEvents(a *Assets, o *gm.Object, makes string) {
	tl := &gm.Timeline{Name: "faru " + makes, Moments: map[int]gm.Handler{}}
	for _, at := range []int{30, 60, 90} {
		tl.Moments[at] = func(r *gm.Room, self *gm.Instance) {
			r.Create(a.Object(makes), self.X, self.Y)
			if at == 90 {
				self.Vars["timelineaktiva"] = 0
			}
		}
	}
	o.Events.Create = func(r *gm.Room, self *gm.Instance) { self.Vars["timelineaktiva"] = 0 }
	o.Events.Step = func(r *gm.Room, self *gm.Instance) {
		if len(r.All(makes)) == 0 && self.Vars["timelineaktiva"] == 0 {
			self.Timeline, self.TimelinePos = tl, 0 // action_set_timeline(…, 0)
			self.Vars["timelineaktiva"] = 1
		}
	}
}

// characterBonus adds the character's collision with b_monster: after the
// netusxu check, stomp for 10, otherwise back to the saved level (no life
// lost).
func characterBonus(a *Assets, o *gm.Object) {
	o.Events.Collision["b_monster"] = func(r *gm.Room, self, other *gm.Instance) {
		if invulnerable(self) {
			return
		}
		if stomp(r, self, other, a, 10) {
			return
		}
		r.G.Play("killcharacter")
		r.SleepMS(1000, roomSpeed)
		r.Goto = "@saved"
	}
}
