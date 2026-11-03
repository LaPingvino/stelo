package game

import (
	"image"
	"image/color"
)

// Text is a line of text drawn over the frame at the screen's real resolution
// (docs/UI.md §1–2: sharp at any size, in Noto Sans), placed in frame
// coordinates (640×480): (X, Y) is the top-left — or the top-centre / top-right
// for Align centre / right — and Size the font size in frame pixels.
type Text struct {
	S      string
	X, Y   float64
	Size   float64
	Weight Weight
	Color  color.Color
	Align  Align
}

// Weight is the Noto Sans cut: regular, bold, or black (titles).
type Weight int

const (
	Regular Weight = iota
	Bold
	Black
)

// Align is the text's horizontal anchor at X.
type Align int

const (
	Left Align = iota
	Centre
	Right
)

// Overlay collects the frame's texts while it is drawn; the program draws
// them on top once the frame is on the screen. A nil *Overlay means no
// overlay (tests, the render tool): the old in-frame lettering is used.
type Overlay struct {
	Rects []Rect // drawn first
	Pics  []Pic  // then pictures (a level's thumbnail …)
	Texts []Text
}

// Pic is a picture drawn into the frame area X, Y, W, H (frame coordinates),
// unsmoothed.
type Pic struct {
	Img        *image.RGBA
	X, Y, W, H float64
}

// Rect is a filled rectangle in frame coordinates (menu highlights, panels),
// drawn at the screen's resolution under the texts.
type Rect struct {
	X, Y, W, H float64
	Color      color.Color
}

// AddRect queues a rectangle.
func (o *Overlay) AddRect(r Rect) { o.Rects = append(o.Rects, r) }

// Add queues a text.
func (o *Overlay) Add(t Text) { o.Texts = append(o.Texts, t) }

// Reset empties it for the next frame.
func (o *Overlay) Reset() { o.Texts, o.Rects, o.Pics = o.Texts[:0], o.Rects[:0], o.Pics[:0] }
