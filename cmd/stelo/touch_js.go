//go:build js

package main

import (
	"strings"
	"syscall/js"
)

// showPadsFromURL: "?pads" in the page's address shows the touch pads from the
// start (for checking the layout without a touch screen).
func showPadsFromURL() bool {
	return strings.Contains(js.Global().Get("location").Get("search").String(), "pads")
}

// pollPhoneKeys takes what the page queued in steloTyped (strings: one
// character each, "\b" Backspace, "\r" Enter) and empties it.
func pollPhoneKeys() phoneKeys {
	var p phoneKeys
	q := js.Global().Get("steloTyped")
	if q.IsUndefined() || q.Length() == 0 {
		return p
	}
	for i := 0; i < q.Length(); i++ {
		s := q.Index(i).String()
		if s == "\r" {
			p.enter = true
		} else {
			p.chars = append(p.chars, []rune(s)...)
		}
	}
	q.Call("splice", 0, q.Length())
	return p
}

// editFromURL: "?edit" in the page's address opens the level editor.
func editFromURL() bool {
	return strings.Contains(js.Global().Get("location").Get("search").String(), "edit")
}

// originalArtFromURL: "?original" in the page's address starts with the 2006
// art instead of the new art (F3 switches at any time).
func originalArtFromURL() bool {
	return strings.Contains(js.Global().Get("location").Get("search").String(), "original")
}

// langFromURL: "?lang=en" in the page's address draws the English text.
func langFromURL() string {
	if strings.Contains(js.Global().Get("location").Get("search").String(), "lang=en") {
		return "en"
	}
	return ""
}
