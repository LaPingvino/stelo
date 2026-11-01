// Package store keeps the highscore table between games: in the browser's
// localStorage on the web, in a JSON file in the user's config directory on
// the desktop (the original kept it in its own data).
package store

import (
	"encoding/json"

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
