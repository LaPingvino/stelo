//go:build !js

package main

import "os"

func showPadsFromURL() bool { return false }

func pollPhoneKeys() phoneKeys { return phoneKeys{} }

func editFromURL() bool { return false }

// originalArtFromURL: on the desktop, STELO_ART=original starts with the 2006 art.
func originalArtFromURL() bool { return os.Getenv("STELO_ART") == "original" }

// langFromURL: on the desktop, STELO_LANG=en draws the English text.
func langFromURL() string { return os.Getenv("STELO_LANG") }
