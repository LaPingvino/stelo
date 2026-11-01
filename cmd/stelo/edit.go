package main

import (
	"io/fs"
	"log"
	"sort"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/editor"
	"github.com/LaPingvino/stelo/internal/game"
	"github.com/LaPingvino/stelo/internal/level"
	"github.com/LaPingvino/stelo/internal/pack"
	"github.com/LaPingvino/stelo/internal/store"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// The level editor (LOOP item 10): F2 opens it on the room being played (the
// level as its file has it), P plays the edited level from its start, F2
// goes back to editing. The mouse (right button erases) or a finger paints;
// the arrow keys scroll; Ctrl+Z undoes; Ctrl+S saves (the web also downloads
// the level file); a level file dropped on the window opens in the editor, a
// dropped pack folder plays; Ctrl+Shift+S saves the levels edited this session
// as a pack.

// editState keeps the editor between plays: one document per level.
type editState struct {
	on     bool
	screen *game.EditScreen
	docs   map[string]*editor.Doc
	down   bool // a press is going on
	// picked receives a level file's text from the browser's file picker
	// (read on the browser's side; taken in Update).
	picked chan string
}

// openEditor edits the named level (keeping earlier edits of it).
func (a *app) openEditor(name string) {
	if a.edit.docs == nil {
		a.edit.docs = map[string]*editor.Doc{}
	}
	if a.edit.picked == nil {
		a.edit.picked = make(chan string, 1)
	}
	d := a.edit.docs[name]
	if d == nil {
		// The saved version (Ctrl+S) if there is one, else the level's file.
		text, ok := store.LoadLevel(name)
		if !ok {
			data, err := stelo.Files.ReadFile("levels/" + name + ".txt")
			if err != nil {
				log.Printf("editor: %v", err)
				return
			}
			text = string(data)
		}
		l, err := level.Parse(text)
		if err != nil {
			log.Printf("editor: %v", err)
			return
		}
		d = editor.New(l)
		a.edit.docs[name] = d
	}
	tool := "block"
	if a.edit.screen != nil {
		tool = a.edit.screen.Tool
	}
	a.edit.screen = game.NewEditScreen(a.assets, d)
	a.edit.screen.Tool = tool
	a.edit.on = true
}

// playEdited plays the level being edited.
func (a *app) playEdited() {
	g, err := game.NewGameWithLevel(stelo.Files, a.assets, a.edit.screen.Doc.L)
	if err != nil {
		log.Printf("editor play: %v", err)
		return
	}
	g.G.Sound = a.game.G.Sound
	g.Highscores, g.SaveHighscores = a.game.Highscores, a.game.SaveHighscores
	a.game = g
	a.edit.on = false
}

// updateEditor handles the editor's input for a frame.
func (a *app) updateEditor() {
	select {
	case text := <-a.edit.picked:
		a.openText(text)
	default:
	}
	e := a.edit.screen
	x, y := ebiten.CursorPosition()
	left := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	right := ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight)
	if ids := ebiten.AppendTouchIDs(nil); len(ids) > 0 { // a finger paints
		x, y = ebiten.TouchPosition(ids[0])
		left = true
	}
	e.MouseX, e.MouseY = x, y
	switch {
	case (left || right) && !a.edit.down:
		a.edit.down = true
		switch e.Press(x, y, right) {
		case "open":
			store.OpenFile(func(name, text string) {
				select { // never block the browser's callback
				case a.edit.picked <- text:
				default:
				}
			})
		case "save":
			a.saveEdited()
		case "play":
			a.playEdited()
			return
		}
	case left || right:
		e.Drag(x, y, right)
	case a.edit.down:
		a.edit.down = false
		e.Release()
	}
	ctrl := ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyMeta)
	switch {
	case ctrl && inpututil.IsKeyJustPressed(ebiten.KeyZ):
		e.Undo()
	case ctrl && ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyS):
		a.savePack()
	case ctrl && inpututil.IsKeyJustPressed(ebiten.KeyS):
		a.saveEdited()
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft):
		e.Scroll(-1, 0)
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowRight):
		e.Scroll(1, 0)
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowUp):
		e.Scroll(0, -1)
	case inpututil.IsKeyJustPressed(ebiten.KeyArrowDown):
		e.Scroll(0, 1)
	case inpututil.IsKeyJustPressed(ebiten.KeyP):
		a.playEdited()
	}
}

// saveEdited keeps the level being edited; on the web it also offers the
// file.
func (a *app) saveEdited() {
	e := a.edit.screen
	name, text := e.Doc.L.Name, e.Doc.Text()
	store.SaveLevel(name, text)
	store.Download(name, text)
}

// openDropped opens a level file dropped on the window (desktop or browser)
// in the editor, replacing any edits of a level with the same name.
func (a *app) openDropped(files fs.FS) {
	// A dropped folder holding a pack plays it.
	if entries, err := fs.ReadDir(files, "."); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			if p, err := pack.Load(files, e.Name()); err == nil {
				a.playPack(p)
				return
			}
		}
	}
	l, err := editor.FindLevel(files)
	if err != nil {
		log.Printf("editor: %v", err)
		return
	}
	if a.edit.docs == nil {
		a.edit.docs = map[string]*editor.Doc{}
	}
	a.edit.docs[l.Name] = editor.New(l)
	a.openEditor(l.Name)
}

// openText opens a level file's text in the editor.
func (a *app) openText(text string) {
	l, err := level.Parse(text)
	if err != nil {
		log.Printf("editor: %v", err)
		return
	}
	if a.edit.docs == nil {
		a.edit.docs = map[string]*editor.Doc{}
	}
	a.edit.docs[l.Name] = editor.New(l)
	a.openEditor(l.Name)
}

// playPack plays a pack from its first level (leaving the editor).
func (a *app) playPack(p *pack.Pack) {
	g, err := game.NewGameWithPack(stelo.Files, a.assets, p.Levels)
	if err != nil {
		log.Printf("pack %s: %v", p.Name, err)
		return
	}
	if a.game != nil {
		g.G.Sound = a.game.G.Sound
		g.Highscores, g.SaveHighscores = a.game.Highscores, a.game.SaveHighscores
	}
	a.game = g
	a.edit.on = false
	ebiten.SetWindowTitle("Stelo — " + p.Name)
}

// savePack saves every level edited this session as a pack (levels in name
// order): a zip download on the web, a pack folder in the config directory
// on the desktop.
func (a *app) savePack() {
	names := make([]string, 0, len(a.edit.docs))
	for n := range a.edit.docs {
		names = append(names, n)
	}
	sort.Strings(names)
	p := &pack.Pack{Name: "Mia Stelo-pako"}
	for _, n := range names {
		p.Levels = append(p.Levels, a.edit.docs[n].L)
	}
	if len(p.Levels) == 0 {
		return
	}
	data, err := p.Zip()
	if err != nil {
		log.Printf("pack: %v", err)
		return
	}
	store.SavePack(p.Name, data)
}
