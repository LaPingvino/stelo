// Package stelo holds the game's data: the recovered original's sprites,
// backgrounds and project description, and the level files.
package stelo

import "embed"

// Files holds original/project.json, original/sprites/*.png,
// original/backgrounds/*.png and levels/*.txt.
//
//go:embed original/project.json original/sprites/*.png original/backgrounds/*.png levels/*.txt
var Files embed.FS
