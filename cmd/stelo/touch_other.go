//go:build !js

package main

import "os"

func showPadsFromURL() bool { return false }

func pollPhoneKeys() phoneKeys { return phoneKeys{} }

func editFromURL() bool { return false }

func novaFromURL() bool { return false }

// langFromURL: on the desktop, STELO_LANG=en draws the English text.
func langFromURL() string { return os.Getenv("STELO_LANG") }
