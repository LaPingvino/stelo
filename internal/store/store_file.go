//go:build !js

package store

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func levelPath(name string) string {
	d := Dir()
	if d == "" {
		return ""
	}
	return filepath.Join(d, "levels", filepath.Base(name)+".txt")
}

func readLevel(name string) (string, bool) {
	p := levelPath(name)
	if p == "" {
		return "", false
	}
	data, err := os.ReadFile(p)
	return string(data), err == nil
}

func writeLevel(name, text string) {
	p := levelPath(name)
	if p == "" || os.MkdirAll(filepath.Dir(p), 0o755) != nil {
		return
	}
	_ = os.WriteFile(p, []byte(text), 0o644)
}

// Download is the web's file offer; on the desktop the level is saved in
// the config directory already (SaveLevel), so nothing more is done.
func Download(name, text string) bool { return false }

// OpenFile is the web's file picker; the desktop takes dropped files instead.
func OpenFile(got func(name, text string)) bool { return false }

// SavePack unpacks a pack's zip into the config directory's packs folder
// (<config>/stelo/packs/<pack>/), where it can be played from.
func SavePack(name string, zipData []byte) {
	d := Dir()
	if d == "" {
		return
	}
	zr, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return
	}
	for _, f := range zr.File {
		if strings.Contains(f.Name, "..") {
			continue
		}
		p := filepath.Join(d, "packs", filepath.FromSlash(f.Name))
		if os.MkdirAll(filepath.Dir(p), 0o755) != nil {
			return
		}
		r, err := f.Open()
		if err != nil {
			return
		}
		data, err := io.ReadAll(r)
		r.Close()
		if err != nil {
			return
		}
		_ = os.WriteFile(p, data, 0o644)
	}
}
