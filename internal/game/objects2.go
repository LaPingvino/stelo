package game

import "github.com/LaPingvino/stelo/internal/gm"

// The player's collisions with the level's things, and the monsters
// (original/objects/character.gml, monster.gml, flyer.gml, monster_dead.gml).

func init() {
	objectEvents["monster"] = monsterEvents
	objectEvents["flyer"] = flyerEvents
	objectEvents["monster_dead"] = monsterDeadEvents
	objectEvents["transirejo"] = func(a *Assets, o *gm.Object) {
		o.Events.Create = func(r *gm.Room, self *gm.Instance) { self.Vars["used"] = 0 }
	}
}

// invulnerable is the cheat netusxu (character's own variable).
func invulnerable(self *gm.Instance) bool { return self.Vars["netusxu"] == 1 }

// die is the character's "destroy the character" branch: sound, 1 s pause,
// a life less, back to the start point.
func die(r *gm.Room, self *gm.Instance) {
	r.G.Play("killcharacter")
	r.SleepMS(1000, roomSpeed)
	r.G.AddLives(-1)
	self.X, self.Y = self.XStart, self.YStart
}

// stomp is "if falling and above the monster's top + 8, squash it".
func stomp(r *gm.Room, self, other *gm.Instance, a *Assets, points int) bool {
	if !(self.VSpeed > 0 && self.Y < other.Y+8) {
		return false
	}
	r.G.Play("killmonster")
	r.G.Score += points
	other.HSpeed, other.VSpeed = 0, 0 // action_move('000010000', 0): stop
	r.Change(other, a.Object("monster_dead"), true)
	return true
}

// characterThings adds the character's collisions with stars, hearts,
// checkpoints, fire, monsters and flyers. Each first checks netusxu.
func characterThings(a *Assets, o *gm.Object) {
	c := o.Events.Collision
	c["stelo_poentoj"] = func(r *gm.Room, self, other *gm.Instance) {
		if invulnerable(self) {
			return
		}
		r.G.Play("getmushroom")
		r.G.Score += 10
		other.Destroyed = true
	}
	c["bonus"] = func(r *gm.Room, self, other *gm.Instance) {
		if invulnerable(self) {
			return
		}
		r.G.Play("finishlevel")
		r.G.AddLives(1)
		other.Destroyed = true
	}
	c["transirejo"] = func(r *gm.Room, self, other *gm.Instance) {
		if invulnerable(self) || other.Vars["used"] != 0 {
			return
		}
		self.XStart, self.YStart = other.X, other.Y
		other.Vars["used"] = 1
		other.Sprite, other.ImageIndex = "transirejo_a", 0
	}
	c["death"] = func(r *gm.Room, self, other *gm.Instance) {
		if invulnerable(self) {
			return
		}
		die(r, self)
	}
	c["monster"] = func(r *gm.Room, self, other *gm.Instance) {
		if invulnerable(self) {
			return
		}
		if !stomp(r, self, other, a, 50) {
			die(r, self)
		}
	}
	c["flyer"] = func(r *gm.Room, self, other *gm.Instance) {
		if invulnerable(self) {
			return
		}
		// In the bonus level a flyer sends you back to the saved level
		// (item 7); elsewhere it costs a life.
		if r.Name == "bonusnivelo" {
			r.G.Play("killcharacter")
			r.SleepMS(1000, roomSpeed)
			r.Goto = "@saved"
			return
		}
		die(r, self)
	}
}

// monsterEvents is monster.gml: walks right at 4; turns at markers (stepping
// back first) and blocks; its sprite follows its direction.
func monsterEvents(a *Assets, o *gm.Object) {
	o.Events.Create = func(r *gm.Room, self *gm.Instance) { self.HSpeed, self.VSpeed = 4, 0 }
	o.Events.EndStep = func(r *gm.Room, self *gm.Instance) {
		if self.HSpeed > 0 {
			setSprite(self, "sprite_monsterr")
		} else {
			setSprite(self, "sprite_monsterl")
		}
	}
	o.Events.Collision = map[string]gm.CollisionHandler{
		"marker": func(r *gm.Room, self, other *gm.Instance) {
			self.X, self.Y = self.XPrevious, self.YPrevious
			self.HSpeed = -self.HSpeed
		},
		"block": func(r *gm.Room, self, other *gm.Instance) { self.HSpeed = -self.HSpeed },
	}
}

// flyerEvents is flyer.gml: flies right at 3; turns at markers and blocks
// (without stepping back on a marker).
func flyerEvents(a *Assets, o *gm.Object) {
	o.Events.Create = func(r *gm.Room, self *gm.Instance) { self.HSpeed, self.VSpeed = 3, 0 }
	o.Events.EndStep = func(r *gm.Room, self *gm.Instance) {
		if self.HSpeed > 0 {
			setSprite(self, "sprite_flyingr")
		} else {
			setSprite(self, "sprite_flyingl")
		}
	}
	turn := func(r *gm.Room, self, other *gm.Instance) { self.HSpeed = -self.HSpeed }
	o.Events.Collision = map[string]gm.CollisionHandler{"marker": turn, "block": turn}
}

// monsterDeadEvents is monster_dead.gml: gone after 10 steps.
func monsterDeadEvents(a *Assets, o *gm.Object) {
	o.Events.Create = func(r *gm.Room, self *gm.Instance) { self.Alarm[0] = 10 }
	o.Events.Alarm[0] = func(r *gm.Room, self *gm.Instance) { self.Destroyed = true }
}
