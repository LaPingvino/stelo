package game

import "image/color"

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
type Overlay struct{ Texts []Text }

// Add queues a text.
func (o *Overlay) Add(t Text) { o.Texts = append(o.Texts, t) }

// Reset empties it for the next frame.
func (o *Overlay) Reset() { o.Texts = o.Texts[:0] }
