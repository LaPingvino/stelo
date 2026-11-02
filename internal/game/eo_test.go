package game

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

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

// Every translated string is one the game really draws (no stale keys): the
// levels' text labels, the HUD and the highscore prompts; and every text
// label in the levels has an English translation.
func TestEnglishCoversTheText(t *testing.T) {
	used := map[string]bool{"Poentoj: ": true, hsGood: true, hsSorry_: true}
	files, _ := filepath.Glob("../../levels/*.txt")
	re := regexp.MustCompile(`^text \d+ \d+ "(.*)"$`)
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(data), "\n") {
			if m := re.FindStringSubmatch(line); m != nil {
				used[m[1]] = true
				if _, ok := english[m[1]]; !ok {
					t.Errorf("%s: no English for %q", f, m[1])
				}
			}
		}
	}
	for k := range english {
		if !used[k] {
			t.Errorf("stale translation key %q", k)
		}
	}
	Lang = "en"
	defer func() { Lang = "eo" }()
	if Tr("Poentoj: ") != "Score: " || Tr("unknown") != "unknown" {
		t.Error("Tr")
	}
}
