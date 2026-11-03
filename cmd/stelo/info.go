package main

import (
	"image/color"

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
			"Stelo +10 · Monstro: saltu sur ĝin +50, tuŝu ĝin alie: −1 vivo",
			"Fluganto: ne venkebla, evitu ĝin · Fajro: −1 vivo · Koro: +1 vivo",
			"Transirejo: tuŝu ĝin, kaj de nun vi rekomencas ĉi tie",
			"Diamanto: al la sekreta bonusnivelo · Elirejo: al la sekva nivelo",
			"Fali el la nivelo ne kostas vivon: vi rekomencas ĉe la transirejo.",
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
			"Star +10 · Monster: jump on it +50, touch it otherwise: −1 life",
			"Flyer: can't be beaten, avoid it · Fire: −1 life · Heart: +1 life",
			"Checkpoint: touch it, and from now on you start again here",
			"Diamond: to the secret bonus level · Exit: to the next level",
			"Falling out of a level costs no life: you start at the checkpoint.",
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
func (a *app) drawInfo(o *game.Overlay, page string) {
	lang := game.Lang
	if lang != "en" {
		lang = "eo"
	}
	lines := infoPages[page][lang]
	o.AddRect(game.Rect{X: 30, Y: 40, W: game.ScreenW - 60, H: game.ScreenH - 80, Color: color.RGBA{0x2a, 0x1c, 0x12, 235}})
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
	o.Add(game.Text{S: game.Tr("Enter / Esc: reen"), X: game.ScreenW / 2, Y: game.ScreenH - 66, Size: 13, Weight: game.Bold, Color: gold, Align: game.Centre})
}
