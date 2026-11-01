package game

import "github.com/LaPingvino/stelo/internal/gm"

// The original objects' events, ported from original/objects/*.gml (see
// docs/STELO-RULES.md). Each function adds an object's handlers to its engine
// object; the comment above each handler quotes the original's action.
var objectEvents = map[string]func(a *Assets, o *gm.Object){}

func init() { objectEvents["character"] = characterEvents }

// roomSpeed is every room's speed (steps per second), for action_sleep.
const roomSpeed = 30

func characterEvents(a *Assets, o *gm.Object) {
	o.Events.Key = map[gm.Key]gm.Handler{
		gm.KeyRight: func(r *gm.Room, self *gm.Instance) { walk(r, self, 4, "sprite_right") },
		gm.KeyLeft:  func(r *gm.Room, self *gm.Instance) { walk(r, self, -4, "sprite_left") },
		// Keyboard <up>: on a ladder, climb 3 if free; otherwise jump if standing.
		gm.KeyUp: func(r *gm.Room, self *gm.Instance) {
			if r.PlaceMeeting(self, self.X, self.Y, "ladder") != nil {
				if r.PlaceFree(self, self.X, self.Y-3) {
					self.Y -= 3
				}
			} else if !r.PlaceFree(self, self.X, self.Y+1) { // action_if_collision(0, 1)
				self.VSpeed = -10
			}
		},
		// Keyboard <down>: on a ladder, climb down 3 if free.
		gm.KeyDown: func(r *gm.Room, self *gm.Instance) {
			if r.PlaceMeeting(self, self.X, self.Y, "ladder") != nil && r.PlaceFree(self, self.X, self.Y+3) {
				self.Y += 3
			}
		},
	}
	o.Events.Step = func(r *gm.Room, self *gm.Instance) {
		// PARTIAL: background_hspeed (the background drifts ±1 by facing) — drawing, later.
		// In the air: gravity 0.5 downwards, else none.
		if r.PlaceFree(self, self.X, self.Y+1) {
			self.Gravity, self.GravityDirection = 0.5, 270
		} else {
			self.Gravity = 0
		}
		// Limit the vertical speed.
		if self.VSpeed > 12 {
			self.VSpeed = 12
		}
		// On a ladder: no gravity, no vertical speed, the climbing sprite.
		if r.PlaceMeeting(self, self.X, self.Y, "ladder") != nil {
			self.Gravity, self.VSpeed = 0, 0
			setSprite(self, "sprite_climbing")
		}
		// The grid snap "if sprite_index == sprite_ladder" never runs: the
		// player never has that sprite (STELO-RULES §3). Cheat codes: item 8.
	}
	o.Events.Collision = map[string]gm.CollisionHandler{
		// Collision with block: move to contact in the direction of motion, max 12; vspeed 0.
		"block": func(r *gm.Room, self, other *gm.Instance) {
			r.MoveContactSolid(self, self.Direction(), 12)
			self.VSpeed = 0
		},
	}
	characterThings(a, o)
	// Outside room: fell below the room → sound, 1 s pause, back to the start
	// point (no life lost). PARTIAL: the sound (item 9).
	o.Events.OutsideRoom = func(r *gm.Room, self *gm.Instance) {
		if self.Y > float64(r.H) {
			r.SleepMS(1000, roomSpeed)
			self.X, self.Y = self.XStart, self.YStart
		}
	}
}

// walk is Keyboard <right>/<left>: face that way; step dx if free, else step
// up to 8 px (dx, −8) and move down to contact, max 8.
func walk(r *gm.Room, self *gm.Instance, dx float64, sprite string) {
	setSprite(self, sprite)
	switch {
	case r.PlaceFree(self, self.X+dx, self.Y):
		self.X += dx
	case r.PlaceFree(self, self.X+dx, self.Y-8):
		self.X += dx
		self.Y -= 8
		r.MoveContactSolid(self, 270, 8)
	}
}

// setSprite is action_set_sprite(sprite, 1): the sprite and image speed 1;
// the frame restarts only when the sprite changes. [engine]
func setSprite(self *gm.Instance, sprite string) {
	if self.Sprite != sprite {
		self.Sprite, self.ImageIndex = sprite, 0
	}
	self.ImageSpeed = 1
}
