//go:build !js

package store

import (
	"os"
	"path/filepath"
)

// Dir is where the desktop game keeps its files (overridable for tests).
var Dir = func() string {
	d, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(d, "stelo")
}

func read() ([]byte, bool) {
	d := Dir()
	if d == "" {
		return nil, false
	}
	data, err := os.ReadFile(filepath.Join(d, key+".json"))
	return data, err == nil
}

func write(data []byte) {
	d := Dir()
	if d == "" || os.MkdirAll(d, 0o755) != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(d, key+".json"), data, 0o644)
}
