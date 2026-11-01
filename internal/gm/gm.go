// Package gm is the part of Game Maker 5.3's engine that Stelo's objects rely
// on: instances with speed and gravity, bounding-box collisions, place_free,
// move_contact_solid, and the order of a step. docs/STELO-RULES.md §2 marks
// these as [engine]; they come from Game Maker's documented behaviour.
//
// ASSUMPTION (to check against the original running): a bounding box is taken
// at the instance's position rounded to the nearest pixel.
package gm

import "math"

// Box is a bounding box in pixels, all edges inclusive (Game Maker's
// bbox_left, bbox_right, bbox_top, bbox_bottom).
type Box struct{ L, R, T, B int }

// Overlaps reports whether two boxes share a pixel.
func (a Box) Overlaps(b Box) bool { return a.L <= b.R && b.L <= a.R && a.T <= b.B && b.T <= a.B }

// Object is what the engine needs of an object: its name, whether it is solid
// and its sprite's bounding box (relative to the instance's position; every
// Stelo sprite has its origin at 0, 0).
type Object struct {
	Name  string
	Solid bool
	BBox  Box
	// Mask reports whether the object collides at all (an object without a
	// sprite has no box).
	Mask bool
	// Events are the object's event handlers.
	Events Events
}

// Instance is one instance in the room.
type Instance struct {
	Obj                  *Object
	X, Y                 float64
	XPrevious, YPrevious float64
	XStart, YStart       float64
	HSpeed, VSpeed       float64
	Gravity              float64
	GravityDirection     float64 // degrees, 270 = down
	Alarm                [12]int // -1 = off
	Destroyed            bool
}

// NewInstance places an instance; its start and previous position are where it
// is placed.
func NewInstance(obj *Object, x, y float64) *Instance {
	in := &Instance{Obj: obj, X: x, Y: y, XPrevious: x, YPrevious: y, XStart: x, YStart: y, GravityDirection: 270}
	for i := range in.Alarm {
		in.Alarm[i] = -1
	}
	return in
}

// BoxAt is the instance's bounding box if it stood at (x, y).
func (in *Instance) BoxAt(x, y float64) Box {
	ix, iy := int(math.Round(x)), int(math.Round(y))
	b := in.Obj.BBox
	return Box{ix + b.L, ix + b.R, iy + b.T, iy + b.B}
}

// Room holds the instances.
type Room struct {
	W, H      int
	Instances []*Instance
}

func (r *Room) live() []*Instance {
	out := r.Instances[:0:0]
	for _, o := range r.Instances {
		if !o.Destroyed {
			out = append(out, o)
		}
	}
	return out
}

// PlaceFree is place_free(x, y): in would overlap no solid instance there.
func (r *Room) PlaceFree(in *Instance, x, y float64) bool {
	if !in.Obj.Mask {
		return true
	}
	b := in.BoxAt(x, y)
	for _, o := range r.live() {
		if o != in && o.Obj.Solid && o.Obj.Mask && b.Overlaps(o.BoxAt(o.X, o.Y)) {
			return false
		}
	}
	return true
}

// PlaceMeeting is place_meeting(x, y, obj): the first instance of the named
// object in overlaps at (x, y), or nil.
func (r *Room) PlaceMeeting(in *Instance, x, y float64, obj string) *Instance {
	if !in.Obj.Mask {
		return nil
	}
	b := in.BoxAt(x, y)
	for _, o := range r.live() {
		if o != in && o.Obj.Name == obj && o.Obj.Mask && b.Overlaps(o.BoxAt(o.X, o.Y)) {
			return o
		}
	}
	return nil
}

// MoveContactSolid is move_contact_solid(dir, maxdist): move one pixel at a
// time in the direction (degrees, 0 = right, 90 = up) until the next pixel
// would touch a solid, at most maxdist pixels. If the instance already touches
// a solid it doesn't move.
func (r *Room) MoveContactSolid(in *Instance, dir, maxdist float64) {
	if !r.PlaceFree(in, in.X, in.Y) {
		return
	}
	dx, dy := math.Cos(dir*math.Pi/180), -math.Sin(dir*math.Pi/180)
	for i := 0; i < int(maxdist); i++ {
		nx, ny := in.X+dx, in.Y+dy
		if !r.PlaceFree(in, nx, ny) {
			return
		}
		in.X, in.Y = nx, ny
	}
}

// Move is the movement part of a step: gravity is added to the speed (in its
// direction), then the speed to the position.
func (in *Instance) Move() {
	in.HSpeed += in.Gravity * math.Cos(in.GravityDirection*math.Pi/180)
	in.VSpeed -= in.Gravity * math.Sin(in.GravityDirection*math.Pi/180)
	in.HSpeed, in.VSpeed = clean(in.HSpeed), clean(in.VSpeed)
	in.X += in.HSpeed
	in.Y += in.VSpeed
}

// clean removes the floating-point noise of cos/sin at right angles.
func clean(v float64) float64 { return math.Round(v*1e9) / 1e9 }

// Direction is the direction of motion in degrees (Game Maker's direction).
func (in *Instance) Direction() float64 {
	if in.HSpeed == 0 && in.VSpeed == 0 {
		return 0
	}
	d := math.Atan2(-in.VSpeed, in.HSpeed) * 180 / math.Pi
	if d < 0 {
		d += 360
	}
	return d
}
