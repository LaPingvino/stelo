package game

import "strings"

// hatReplacer turns the x-system the 2006 texts are written in into the true
// Esperanto letters (cx → ĉ, …; Joop, 2026-11-01). Only displayed text goes
// through it: the recovered original and the level files keep the x-system,
// and the cheat codes are still matched as typed.
var hatReplacer = strings.NewReplacer(
	"cx", "ĉ", "gx", "ĝ", "hx", "ĥ", "jx", "ĵ", "sx", "ŝ", "ux", "ŭ",
	"Cx", "Ĉ", "Gx", "Ĝ", "Hx", "Ĥ", "Jx", "Ĵ", "Sx", "Ŝ", "Ux", "Ŭ",
	"CX", "Ĉ", "GX", "Ĝ", "HX", "Ĥ", "JX", "Ĵ", "SX", "Ŝ", "UX", "Ŭ",
)

// hats converts x-system Esperanto to the true letters.
func hats(s string) string { return hatReplacer.Replace(s) }

// Hats is hats for other packages: typed x-system letters (cx, gx, …) as
// the true Esperanto letters (pack names and authors are kept that way).
func Hats(s string) string { return hats(s) }
