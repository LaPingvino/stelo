// Package store keeps the highscore table between games: in the browser's
// localStorage on the web, in a JSON file in the user's config directory on
// the desktop (the original kept it in its own data).
package store

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/LaPingvino/stelo/internal/game"
)

const key = "stelo-highscores"

// LoadHighscores returns the saved table, or the cleared one.
func LoadHighscores() game.Highscores {
	h := game.ClearedHighscores()
	if data, ok := read(); ok {
		var saved game.Highscores
		if json.Unmarshal(data, &saved) == nil {
			h = saved
		}
	}
	return h
}

// SaveHighscores stores the table.
func SaveHighscores(h game.Highscores) {
	if data, err := json.Marshal(h); err == nil {
		write(data)
	}
}

// SaveLevel keeps an edited level (the level file's text) under its name;
// LoadLevel gives it back. The editor saves with Ctrl+S and opens a level's
// saved version when there is one.
func SaveLevel(name, text string) { writeLevel(name, text) }

// LoadLevel returns a saved level's text.
func LoadLevel(name string) (string, bool) { return readLevel(name) }

// Progress is the furthest level reached in a pack (1-based; 0 none).
func Progress(pack string) int {
	data, ok := readKey("progress-" + pack)
	if !ok {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return n
}

// Reached records a level reached in a pack, keeping the furthest.
func Reached(pack string, n int) {
	if n > Progress(pack) {
		writeKey("progress-"+pack, []byte(strconv.Itoa(n)))
	}
}

// Muted is the sound button's setting (off: false).
func Muted() bool {
	data, ok := readKey("muted")
	return ok && strings.TrimSpace(string(data)) == "1"
}

// SetMuted keeps the sound button's setting.
func SetMuted(m bool) {
	v := "0"
	if m {
		v = "1"
	}
	writeKey("muted", []byte(v))
}
