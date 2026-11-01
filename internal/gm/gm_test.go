package gm

import "testing"

var (
	block  = &Object{Name: "block", Solid: true, Mask: true, BBox: Box{0, 31, 0, 31}}
	player = &Object{Name: "character", Mask: true, BBox: Box{7, 23, 2, 31}} // sprite_right's box
	ladder = &Object{Name: "ladder", Mask: true, BBox: Box{6, 25, 0, 31}}
)

func TestBoxesInclusive(t *testing.T) {
	a := Box{0, 31, 0, 31}
	if !a.Overlaps(Box{31, 40, 0, 5}) || a.Overlaps(Box{32, 40, 0, 5}) {
		t.Error("edges are inclusive: touching at pixel 31 overlaps, 32 doesn't")
	}
}

func TestPlaceFreeAndMeeting(t *testing.T) {
	r := &Room{W: 640, H: 480}
	floor := NewInstance(block, 64, 96)
	p := NewInstance(player, 64, 64)
	l := NewInstance(ladder, 64, 64)
	r.Instances = []*Instance{floor, p, l}
	if !r.PlaceFree(p, 64, 64) {
		t.Error("standing on the floor is free (box bottom 95, floor top 96)")
	}
	if r.PlaceFree(p, 64, 65) {
		t.Error("one pixel down overlaps the floor")
	}
	if r.PlaceMeeting(p, 64, 64, "ladder") != l {
		t.Error("meets the ladder")
	}
	if r.PlaceMeeting(p, 64+40, 64, "ladder") != nil {
		t.Error("40 px right is off the ladder")
	}
}

func TestMoveContactSolid(t *testing.T) {
	r := &Room{W: 640, H: 480}
	floor := NewInstance(block, 64, 96)
	p := NewInstance(player, 64, 50)
	r.Instances = []*Instance{floor, p}
	r.MoveContactSolid(p, 270, 12)
	if p.Y != 62 {
		t.Errorf("12 px down from 50: y = %v, want 62 (not yet touching)", p.Y)
	}
	r.MoveContactSolid(p, 270, 12)
	if p.Y != 64 {
		t.Errorf("stops touching the floor: y = %v, want 64", p.Y)
	}
}

func TestGravityAndCap(t *testing.T) {
	p := NewInstance(player, 0, 0)
	p.Gravity = 0.5
	for range 4 {
		p.Move()
	}
	// vspeed 0.5, 1, 1.5, 2 → y 0.5 + 1 + 1.5 + 2 = 5
	if p.VSpeed != 2 || p.Y != 5 || p.X != 0 {
		t.Errorf("after 4 steps: vspeed %v, y %v, x %v", p.VSpeed, p.Y, p.X)
	}
	if d := p.Direction(); d != 270 {
		t.Errorf("falling direction %v, want 270", d)
	}
}
