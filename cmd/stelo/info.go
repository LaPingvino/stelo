package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/LaPingvino/stelo/internal/game"
)

// The title's "Kiel ludi" (how to play) and "Pri" (about) pages — the texts of
// docs/TEXTS.md (UI step U3), in Esperanto and English. DRAFT texts, put in
// so Joop can correct them in place (his choice, 2026-11-04).
var infoPages = map[string]map[string][]string{
	"help": {
		"eo": {
			"# Kiel ludi",
			"Konduku Stelon tra ĉiu nivelo ĝis la elirejo. Kolektu stelojn",
			"por poentoj, evitu fajron kaj flugantojn, kaj saltu sur monstrojn.",
			"",
			"* Klavoj: la sagoj — flanken iri, supren salti aŭ grimpi,",
			"malsupren malsupreniri ŝtupetaron · Esc aŭ P paŭzo · M sono",
			"F11 plena ekrano",
			"* Telefono: la butonoj ĉe la randoj; frapeto = iu ajn klavo.",
			"",
		},
		"en": {
			"# How to play",
			"Lead Stelo through each level to the exit. Collect stars",
			"for points, avoid fire and flyers, and jump on monsters.",
			"",
			"* Keys: the arrows — left / right walk, up jump or climb,",
			"down climb down a ladder · Esc or P pause · M sound",
			"F11 full screen",
			"* Phone: the buttons at the edges; a tap = any key.",
			"",
		},
	},
	"about": {
		"eo": {
			"# Pri",
			"* Stelo — platforma ludo de Joop Kiefte.",
			"Farita en 2006 per Game Maker; refarita en 2026 en Go",
			"(Ebitengine), paŝo post paŝo laŭ la originalo.",
			"",
			"Nova arto: @lucasserapiao. La arto de 2006 restas elektebla.",
			"Muziko kaj sonoj: el la ludo de 2006 (muziko per FluidR3, MIT).",
			"Tiparo: Noto Sans (SIL OFL). Sono malkodita per gopus.",
			"",
			"Fontkodo: github.com/LaPingvino/stelo",
		},
		"en": {
			"# About",
			"* Stelo — a platform game by Joop Kiefte.",
			"Made in 2006 with Game Maker; remade in 2026 in Go",
			"(Ebitengine), step by step after the original.",
			"",
			"New art: @lucasserapiao. The 2006 art can still be chosen.",
			"Music and sounds: from the 2006 game (music via FluidR3, MIT).",
			"Font: Noto Sans (SIL OFL). Sound decoded with gopus.",
			"",
			"Source: github.com/LaPingvino/stelo",
		},
	},
}

// drawInfo draws an info page on a panel over the title: "# " a heading,
// "* " a bold line; Enter, Esc or a click goes back.
func (a *app) drawInfo(o *game.Overlay, page string, frame func(string, int) *image.RGBA) {
	lang := game.Lang
	if lang != "en" {
		lang = "eo"
	}
	lines := infoPages[page][lang]
	// The panel goes into the frame (the overlay's rects would cover the
	// legend's sprites): a dark wood brown at 92%.
	panel := image.Rect(30, 40, game.ScreenW-30, game.ScreenH-40)
	draw.Draw(a.frame, panel, image.NewUniform(color.RGBA{0x27, 0x1a, 0x11, 235}), image.Point{}, draw.Over)
	cream, gold := color.RGBA{0xff, 0xf4, 0xd8, 255}, color.RGBA{0xff, 0xd0, 0x3a, 255}
	y := 58.0
	for _, l := range lines {
		switch {
		case len(l) > 2 && l[:2] == "# ":
			o.Add(game.Text{S: l[2:], X: game.ScreenW / 2, Y: y, Size: 26, Weight: game.Black, Color: gold, Align: game.Centre})
			y += 40
		case len(l) > 2 && l[:2] == "* ":
			o.Add(game.Text{S: l[2:], X: 48, Y: y, Size: 14, Weight: game.Bold, Color: cream})
			y += 22
		case l == "":
			y += 10
		default:
			o.Add(game.Text{S: l, X: 48, Y: y, Size: 14, Color: cream})
			y += 22
		}
	}
	if page == "help" {
		drawLegend(a.frame, o, frame, y, lang)
	}
	o.Add(game.Text{S: game.Tr("Enter / Esc: reen"), X: game.ScreenW / 2, Y: game.ScreenH - 56, Size: 13, Weight: game.Bold, Color: gold, Align: game.Centre})
}

// legend is the help page's legend: each thing with its own sprite (UI.md §4),
// in two columns — the texts of docs/TEXTS.md.
var legend = []struct {
	sprite string
	eo, en string
}{
	{"sprite_stelo_poentoj", "Stelo: +10", "Star: +10"},
	{"sprite_monsterl", "Monstro: saltu sur ĝin +50", "Monster: jump on it +50"},
	{"sprite_flyingl", "Fluganto: evitu ĝin", "Flyer: avoid it"},
	{"sprite_death", "Fajro: −1 vivo", "Fire: −1 life"},
	{"sprite_bonus", "Koro: +1 vivo", "Heart: +1 life"},
	{"transirejo_na", "Transirejo: rekomencu ĉi tie", "Checkpoint: start again here"},
	{"diamanto", "Diamanto: al la bonusnivelo", "Diamond: to the bonus level"},
	{"sprite_level", "Elirejo: al la sekva nivelo", "Exit: to the next level"},
}

// legendFoot is the line under the legend.
var legendFoot = map[string]string{
	"eo": "Fali el la nivelo ne kostas vivon: vi rekomencas ĉe la transirejo.",
	"en": "Falling out of a level costs no life: you start at the checkpoint.",
}

// drawLegend draws the legend's sprites into the frame (frame gives a
// sprite's picture in the chosen art) and its texts into the overlay, from y.
func drawLegend(dst *image.RGBA, o *game.Overlay, frame func(string, int) *image.RGBA, y float64, lang string) {
	cream := color.RGBA{0xff, 0xf4, 0xd8, 255}
	for i, l := range legend {
		x := 48 + (i/4)*280
		row := y + float64(i%4)*36
		if img := frame(l.sprite, 0); img != nil {
			draw.Draw(dst, img.Bounds().Add(image.Pt(x, int(row))), img, img.Bounds().Min, draw.Over)
		}
		t := l.eo
		if lang == "en" {
			t = l.en
		}
		o.Add(game.Text{S: t, X: float64(x + 40), Y: row + 8, Size: 14, Color: cream})
	}
	o.Add(game.Text{S: legendFoot[lang], X: 48, Y: y + 4*36 + 4, Size: 14, Color: cream})
}
