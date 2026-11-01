package gm

import "sort"

// Key is a keyboard key as Game Maker numbers them (virtual key codes).
type Key int

// Keys Stelo and Pacmanio use.
const (
	KeyLeft  Key = 37
	KeyUp    Key = 38
	KeyRight Key = 39
	KeyDown  Key = 40
	KeySpace Key = 32
	KeyPgDn  Key = 34
	KeyEnter Key = 13
)

// Input is the keyboard as a step sees it.
type Input interface {
	Held(Key) bool    // keyboard event: the key is down
	Pressed(Key) bool // key press event: it went down since the last step
}

// Handler is an event's code: the room and the instance it runs as (self).
type Handler func(r *Room, self *Instance)

// CollisionHandler also gets the other instance.
type CollisionHandler func(r *Room, self, other *Instance)

// Events are an object's event handlers.
type Events struct {
	Create      Handler
	Alarm       [12]Handler
	Key         map[Key]Handler // keyboard (held)
	KeyPress    map[Key]Handler
	Step        Handler
	OutsideRoom Handler
	Collision   map[string]CollisionHandler // by the other object's name
	EndStep     Handler
}

// Create places a new instance and runs its create event (instance_create).
func (r *Room) Create(obj *Object, x, y float64) *Instance {
	in := NewInstance(obj, x, y)
	r.Instances = append(r.Instances, in)
	if obj.Events.Create != nil {
		obj.Events.Create(r, in)
	}
	return in
}

// Step runs one step in Game Maker's order [engine]: alarms, keyboard, key
// press, step, movement, outside room, collisions, end step. Each phase runs
// over the instances in creation order; instances created during a phase
// take part from the next phase on, destroyed ones are skipped and removed at
// the end of the step.
//
// ASSUMPTION (from the GM manual as recalled, to check against the original):
// on a collision with a solid instance, the moving instance is first put back
// at its previous position, then the event runs, and afterwards it moves on by
// its (possibly changed) speed if that place is free.
func (r *Room) Step(in Input) {
	each := func(f func(*Instance)) {
		for _, i := range r.Instances[:len(r.Instances):len(r.Instances)] {
			if !i.Destroyed {
				f(i)
			}
		}
	}
	each(func(i *Instance) { // alarms
		for a := range i.Alarm {
			if i.Alarm[a] < 0 {
				continue
			}
			if i.Alarm[a]--; i.Alarm[a] == 0 {
				i.Alarm[a] = -1
				if h := i.Obj.Events.Alarm[a]; h != nil {
					h(r, i)
				}
			}
		}
	})
	each(func(i *Instance) {
		for _, k := range sortedKeys(i.Obj.Events.Key) {
			if in.Held(k) {
				i.Obj.Events.Key[k](r, i)
			}
		}
	})
	each(func(i *Instance) {
		for _, k := range sortedKeys(i.Obj.Events.KeyPress) {
			if in.Pressed(k) {
				i.Obj.Events.KeyPress[k](r, i)
			}
		}
	})
	each(func(i *Instance) {
		if i.Obj.Events.Step != nil {
			i.Obj.Events.Step(r, i)
		}
	})
	each(func(i *Instance) {
		i.XPrevious, i.YPrevious = i.X, i.Y
		i.Move()
	})
	each(func(i *Instance) {
		if i.Obj.Events.OutsideRoom != nil && r.outside(i) {
			i.Obj.Events.OutsideRoom(r, i)
		}
	})
	each(func(i *Instance) {
		for _, name := range sortedNames(i.Obj.Events.Collision) {
			h := i.Obj.Events.Collision[name]
			for _, o := range r.Instances {
				if o == i || o.Destroyed || i.Destroyed || o.Obj.Name != name || !i.Obj.Mask || !o.Obj.Mask {
					continue
				}
				if !i.BoxAt(i.X, i.Y).Overlaps(o.BoxAt(o.X, o.Y)) {
					continue
				}
				if o.Obj.Solid {
					i.X, i.Y = i.XPrevious, i.YPrevious
					h(r, i, o)
					if nx, ny := i.X+i.HSpeed, i.Y+i.VSpeed; (i.HSpeed != 0 || i.VSpeed != 0) && r.PlaceFree(i, nx, ny) {
						i.X, i.Y = nx, ny
					}
					continue
				}
				h(r, i, o)
			}
		}
	})
	each(func(i *Instance) {
		if i.Obj.Events.EndStep != nil {
			i.Obj.Events.EndStep(r, i)
		}
	})
	live := r.Instances[:0]
	for _, i := range r.Instances {
		if !i.Destroyed {
			live = append(live, i)
		}
	}
	r.Instances = live
}

// outside reports whether the instance is entirely outside the room (its box,
// or its position without one).
func (r *Room) outside(i *Instance) bool {
	b := Box{int(i.X), int(i.X), int(i.Y), int(i.Y)}
	if i.Obj.Mask {
		b = i.BoxAt(i.X, i.Y)
	}
	return b.R < 0 || b.L >= r.W || b.B < 0 || b.T >= r.H
}

// Events run in a fixed order: keys by key code, collisions by the other
// object's name. [engine] Game Maker orders them by key code and by object
// index; ASSUMPTION: the name order is close enough until the object index is
// carried over (no Stelo object has two collisions that interact in one step).
func sortedKeys(m map[Key]Handler) []Key {
	ks := make([]Key, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Slice(ks, func(a, b int) bool { return ks[a] < ks[b] })
	return ks
}

func sortedNames(m map[string]CollisionHandler) []string {
	ns := make([]string, 0, len(m))
	for n := range m {
		ns = append(ns, n)
	}
	sort.Strings(ns)
	return ns
}
