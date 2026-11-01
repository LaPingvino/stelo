// Package pack reads level packs (docs/LEVEL-FORMAT.md "Packs"), as in The
// Puzzle Pits: the original levels are the first pack; any other is a folder
// of level files with an optional pack.txt (name= and author= lines). The
// levels play in file-name order.
package pack

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/LaPingvino/stelo/internal/level"
)

// Pack is a named list of levels in play order.
type Pack struct {
	Name, Author string
	Levels       []*level.Level
}

// Load reads the pack in dir of fsys: pack.txt (if any) and every other
// *.txt as a level, sorted by file name. A file that doesn't read as a
// level is an error naming it.
func Load(fsys fs.FS, dir string) (*Pack, error) {
	p := &Pack{Name: path.Base(dir)}
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		n := e.Name()
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(n), ".txt") {
			continue
		}
		if strings.EqualFold(n, "pack.txt") {
			data, err := fs.ReadFile(fsys, path.Join(dir, n))
			if err != nil {
				return nil, err
			}
			p.readInfo(string(data))
			continue
		}
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		data, err := fs.ReadFile(fsys, path.Join(dir, n))
		if err != nil {
			return nil, err
		}
		l, err := level.Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", n, err)
		}
		p.Levels = append(p.Levels, l)
	}
	if len(p.Levels) == 0 {
		return nil, fmt.Errorf("%s: no levels", dir)
	}
	return p, nil
}

// readInfo takes name= and author= from pack.txt.
func (p *Pack) readInfo(s string) {
	for _, line := range strings.Split(s, "\n") {
		k, v, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(k)) {
		case "name":
			p.Name = strings.TrimSpace(v)
		case "author":
			p.Author = strings.TrimSpace(v)
		}
	}
}

// Original is the first pack: the 2006 game's five levels (levels/ in the
// game's files), by Joop.
func Original(fsys fs.FS) (*Pack, error) {
	p := &Pack{Name: "Stelo", Author: "LaPingvino"}
	for i := 1; i <= 5; i++ {
		data, err := fs.ReadFile(fsys, fmt.Sprintf("levels/nivelo%d.txt", i))
		if err != nil {
			return nil, err
		}
		l, err := level.Parse(string(data))
		if err != nil {
			return nil, err
		}
		p.Levels = append(p.Levels, l)
	}
	return p, nil
}
