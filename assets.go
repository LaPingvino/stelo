// Package stelo holds the game's data: the recovered original's sprites,
// backgrounds and project description, and the level files.
package stelo

import "embed"

// Files holds original/project.json, original/sprites/*.png,
// original/backgrounds/*.png, audio/effects/*.opus and audio/music/*.opus (the
// sounds, converted by scripts/convert-audio.sh), levels/*.txt and
// fonts/*.ttf (Liberation Sans,
// SIL Open Font License, standing in for the original's Arial).
//
//go:embed original/project.json original/sprites/*.png original/backgrounds/*.png audio/effects/*.opus audio/music/*.opus levels/*.txt fonts/*.ttf art/nova/*.png
var Files embed.FS
