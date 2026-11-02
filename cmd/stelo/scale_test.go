package main

import "testing"

// fit: whole steps when they fill at least 90% of the room, else the largest
// fraction; the frame centred; toFrame inverts the placement.
func TestFit(t *testing.T) {
	for _, c := range []struct {
		w, h  int
		scale float64
	}{
		{1280, 960, 2},       // exactly 2×
		{1920, 1080, 2},      // 2.25 fits; 2 fills 89% … fraction
		{2532, 1170, 2.4375}, // a phone on its side at 3× density
		{640, 480, 1},
		{1366, 900, 1.875},
	} {
		v := fit(c.w, c.h)
		want := c.scale
		if c.w == 1920 {
			want = 2.25
		}
		if v.scale != want {
			t.Errorf("%d×%d: scale %v, want %v", c.w, c.h, v.scale, want)
		}
		view = v
		if x, y := toFrame(int(v.ox+320*v.scale), int(v.oy+240*v.scale)); x != 320 || y != 240 {
			t.Errorf("%d×%d: centre maps to %d,%d", c.w, c.h, x, y)
		}
	}
}
