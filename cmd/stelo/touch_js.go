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
