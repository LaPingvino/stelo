// Package stelo holds the game's data: the recovered original's sprites,
// backgrounds and project description, and the level files.
package stelo

import "embed"

// Files holds original/project.json, original/sprites/*.png,
// original/backgrounds/*.png, original/sounds/*.wav (the effects), levels/*.txt and
// fonts/*.ttf (Liberation Sans,
// SIL Open Font License, standing in for the original's Arial).
//
//go:embed original/project.json original/sprites/*.png original/backgrounds/*.png original/sounds/*.wav levels/*.txt fonts/*.ttf art/nova/*.png
var Files embed.FS
