package game

import "testing"

func TestHats(t *testing.T) {
	for in, want := range map[string]string{
		"Cxi tiu ludo ne efektivigxus": "Ĉi tiu ludo ne efektiviĝus",
		"Sxtupetaro":                   "Ŝtupetaro",
		"Ho ve... Vi ne suficxe bone":  "Ho ve... Vi ne sufiĉe bone",
		"auxdi, ehxo, jxus, NUN IRU":   "aŭdi, eĥo, ĵus, NUN IRU",
		"FINLUDI LA TUTAN LUDON":       "FINLUDI LA TUTAN LUDON",
	} {
		if got := hats(in); got != want {
			t.Errorf("hats(%q) = %q, want %q", in, got, want)
		}
	}
}
