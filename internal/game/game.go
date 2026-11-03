package game

import (
	"errors"
	"fmt"
	"github.com/LaPingvino/stelo/internal/level"
	"golang.org/x/image/font"
	"image"
	"io/fs"
	"strconv"

	"github.com/LaPingvino/stelo/internal/gm"
)

// RoomOrder is the original's room order (project.json room_order); the
// bonus level is reached only through the diamond (or a cheat).
var RoomOrder = []string{"klarigo", "nivelo1", "nivelo2", "nivelo3", "nivelo4", "nivelo5", "danklisto", "finlisto", "bonusnivelo"}

// Game runs the rooms one after another with the state that lives across
// them: score, lives, and the level saved while the bonus level is played
// (the original makes that room persistent and returns to global.ludnivelo).
type Game struct {
	levelStart levelStart
	FS         fs.FS
	Assets     *Assets
	G          *gm.Globals
	World      *World
	saved      *World // the persistent level left for the bonus level
	Font       font.Face
	TextFont   font.Face
	// Title: the title picture is up (game start); waitRoom: klarigo waits for
	// a key (klarigilo's keyboard_wait).
	Title    bool
	waitRoom bool
	// The highscore table, its screen state and font; SaveHighscores stores
	// the table when it changes (nil: kept in memory only).
	Highscores     Highscores
	hs             hsState
	HSFont         font.Face
	SaveHighscores func(Highscores)
	// Levels are edited levels played instead of the files (by name).
	Levels map[string]*level.Level
	// Order is the rooms' play order when playing a pack (nil: RoomOrder).
	Order []string
	// OnLevel, if set, is told the pack level (1-based) entered when playing
	// a pack or the original's levels (progress per pack).
	OnLevel func(n int)
	// OnRoom, if set, is told every room entered (the title's level pick
	// shows the bonus level once it was found).
	OnRoom func(room string)
	// OnRestart, when set, replaces game_restart: the program goes back to
	// its title screen.
	OnRestart func()
	// Skin is the drawing skin for every room (nil: the original art).
	Skin *Skin
}

// NewGame starts at a room (normally RoomOrder[0]).
// order is the rooms' play order: a pack's (Order) or the original's.
func (g *Game) order() []string {
	if g.Order != nil {
		return g.Order
	}
	return RoomOrder
}

// NewGameWithPack plays a pack: its levels in order (as rooms "pack:1",
// "pack:2", …, so their names can't clash with the original's), then the
// original's credits and highscores, as after level 5.
func NewGameWithPack(fsys fs.FS, a *Assets, levels []*level.Level) (*Game, error) {
	if len(levels) == 0 {
		return nil, errors.New("an empty pack")
	}
	g, err := NewGame(fsys, a, "")
	if g == nil {
		return nil, err
	}
	g.Levels = map[string]*level.Level{}
	for i, l := range levels {
		cp, err := level.Parse(level.Format(l))
		if err != nil {
			return nil, err
		}
		cp.Name = fmt.Sprintf("pack:%d", i+1)
		g.Levels[cp.Name] = cp
		g.Order = append(g.Order, cp.Name)
	}
	g.Order = append(g.Order, "danklisto", "finlisto")
	return g, g.enter(g.Order[0])
}

// NewGameWithLevel plays an edited level (a copy of it) from its start, as
// the editor's Play: the room's own name, the rest of the game as usual.
func NewGameWithLevel(fsys fs.FS, a *Assets, l *level.Level) (*Game, error) {
	cp, err := level.Parse(level.Format(l))
	if err != nil {
		return nil, err
	}
	g, err := NewGame(fsys, a, "")
	if g == nil {
		return nil, err
	}
	g.Title = false
	g.Levels = map[string]*level.Level{cp.Name: cp}
	return g, g.enter(cp.Name)
}

func NewGame(fsys fs.FS, a *Assets, room string) (*Game, error) {
	g := &Game{FS: fsys, Assets: a, G: &gm.Globals{Lives: 3, Vars: map[string]float64{}}}
	if f, err := HUDFont(fsys); err == nil {
		g.Font = f
	}
	if f, err := TextFont(fsys); err == nil {
		g.TextFont = f
	}
	if f, err := HighscoreFont(fsys); err == nil {
		g.HSFont = f
	}
	g.Highscores = ClearedHighscores()
	// klarigilo's game start event (the first room holds the klarigilo):
	// the komenco sound and the title picture until a key.
	if room == RoomOrder[0] {
		g.Title = true
		g.G.Play("@sound:komenco")
	}
	if room == "" {
		return g, nil
	}
	return g, g.enter(room)
}

func (g *Game) enter(room string) error {
	if g.OnRoom != nil {
		g.OnRoom(room)
	}
	if n := g.levelNumber(room); n > 0 {
		g.levelStart = levelStart{n: n, score: g.G.Score, lives: g.G.Lives}
		if g.OnLevel != nil {
			g.OnLevel(n)
		}
	} else {
		g.levelStart = levelStart{}
	}
	var w *World
	if l := g.Levels[room]; l != nil { // an edited level (the editor's Play)
		w = NewWorld(g.Assets, l)
	} else {
		var err error
		if w, err = LoadLevel(g.FS, g.Assets, room); err != nil {
			return err
		}
	}
	w.Room.G = g.G
	w.Skin = g.Skin
	w.Font, w.TextFont = g.Font, g.TextFont
	g.World = w
	switch n := g.levelNumber(room); {
	case n > 0:
		w.Label = Tr("Nivelo") + " " + strconv.Itoa(n)
	case room == "bonusnivelo":
		w.Label = Tr("Bonusnivelo")
	}
	// The persistent life controller exists from the room that places it
	// (nivelo1) until it is destroyed.
	if len(w.Room.All("life_controller")) > 0 {
		g.G.Vars["life_controller_alive"] = 1
	}
	if code := roomCode[room]; code != nil {
		code(g)
	}
	if len(w.Room.All("klarigilo")) > 0 {
		g.waitRoom = true
	}
	if len(w.Room.All("plejbonaj_listo")) > 0 {
		g.startHighscores()
	}
	return nil
}

// roomCode is the rooms' creation code, ported (levels/*.txt "code" lines).
var roomCode = map[string]func(g *Game){
	// danklisto: with (life_controller) {instance_destroy();}
	// sound_stop_all(); sound_loop(LaEspero);
	"danklisto": func(g *Game) {
		g.G.Vars["life_controller_gone"] = 1
		g.G.Play("@music:LaEspero")
	},
}

// AnyKey is the part of the input keyboard_wait needs: a key went down this
// step (any key).
type AnyKey interface{ AnyPressed() bool }

// Typer is the part of the input keyboard_string needs: the characters typed
// since the last step.
type Typer interface{ Typed() []rune }

// Waiting reports a keyboard_wait in progress: the title picture (Title) or
// the rules screen. The game draws but doesn't step while waiting.
func (g *Game) Waiting() bool { return g.Title || g.waitRoom }

// Step runs one step and then any room change it asked for.
func (g *Game) Step(in gm.Input) error {
	if g.World != nil {
		g.World.Steps++
	}
	// klarigilo (original/objects/klarigilo.gml): at game start the title
	// picture stelo_komenco and keyboard_wait(); in klarigo's room start
	// screen_redraw(); keyboard_wait(); room_goto_next().
	if g.Waiting() {
		if a, ok := in.(AnyKey); ok && a.AnyPressed() {
			if g.Title {
				g.Title = false
				g.G.Play("@music:muziko1") // end komenco, loop muziko1
				return nil
			}
			g.waitRoom = false
			g.World.Room.Goto = "@next"
			return g.change()
		}
		return nil
	}
	if t, ok := in.(Typer); ok {
		g.G.Type(t.Typed())
	}
	if g.hs != hsNone {
		anyKey := false
		if a, ok := in.(AnyKey); ok {
			anyKey = a.AnyPressed()
		}
		if g.stepHighscores(in.Pressed(gm.KeyEnter), anyKey) {
			g.restart()
			return nil
		}
		if g.hs != hsTable { // the name prompt and "Ho ve" block the room
			return nil
		}
	}
	g.World.Step(in)
	if g.G.Vars["highscore_clear"] == 1 { // the pasispasinteco cheat
		g.G.Vars["highscore_clear"] = 0
		g.Highscores = ClearedHighscores()
		if g.SaveHighscores != nil {
			g.SaveHighscores(g.Highscores)
		}
	}
	if g.G.NoMoreLives { // life_controller: "no more lives" → finlisto
		g.G.NoMoreLives = false
		g.World.Room.Goto = "finlisto"
	}
	return g.change()
}

// change performs the room change the room asked for, if any.
func (g *Game) change() error {
	to := g.World.Room.Goto
	if to == "" || g.World.Room.Sleep > 0 {
		return nil // a room change waits until the pause is over
	}
	g.World.Room.Goto = ""
	switch to {
	case "@next":
		order := g.order()
		for i, n := range order {
			if n == g.World.Level.Name && i+1 < len(order) {
				return g.enter(order[i+1])
			}
		}
		return nil
	case "@bonus": // the room becomes persistent: keep it as it is
		g.saved = g.World
		return g.enter("bonusnivelo")
	case "@saved":
		if g.saved != nil {
			g.World, g.saved = g.saved, nil
			return nil
		}
		return nil
	}
	return g.enter(to)
}

// characterRooms adds the character's room-changing collisions.
func characterRooms(a *Assets, o *gm.Object) {
	{
		c := o.Events.Collision
		// Collision with levelexit: sound, 1 s pause, the next room (if any).
		c["levelexit"] = func(r *gm.Room, self, other *gm.Instance) {
			if invulnerable(self) {
				return
			}
			r.G.Play("finishlevel")
			r.SleepMS(1000, roomSpeed)
			r.Goto = "@next"
		}
		// Collision with iru_al_bonuso (the diamond): the room is kept as it
		// is, the diamond disappears, the bonus music, the bonus level.
		c["iru_al_bonuso"] = func(r *gm.Room, self, other *gm.Instance) {
			other.Destroyed = true
			r.G.Play("@music:bonussono")
			r.Goto = "@bonus"
		}
	}
}

// restart is game_restart(): a new game from the first room, the highscore
// table kept.
func (g *Game) restart() {
	if g.OnRestart != nil { // the program's title screen takes over (docs/UI.md)
		g.OnRestart()
		return
	}
	hs, save, snd := g.Highscores, g.SaveHighscores, g.G.Sound
	ng, err := NewGame(g.FS, g.Assets, RoomOrder[0])
	if err != nil {
		return
	}
	ng.Highscores, ng.SaveHighscores, ng.G.Sound = hs, save, snd
	ng.OnLevel = g.OnLevel
	ng.Skin = g.Skin
	*g = *ng
}

// Draw draws the current screen: the title, or the room with the highscore
// screen over it when it is up.
func (g *Game) Draw(dst *image.RGBA) {
	if g.Title {
		g.World.DrawTitle(dst)
		return
	}
	g.World.Draw(dst)
	g.drawHighscores(dst)
}

// levelNumber is a room's place among the levels being played (1-based):
// pack:n for a pack, nivelo1–5 for the original; 0 for other rooms.
// EnterLevel starts level n (1-based, the pack's or the original's order),
// leaving any title or waiting room: the new title's Play and Continue.
func (g *Game) EnterLevel(n int) error {
	g.Title, g.waitRoom = false, false
	rooms := RoomOrder
	if g.Order != nil {
		rooms = g.Order
	}
	for _, r := range rooms {
		if g.levelNumber(r) == n {
			return g.enter(r)
		}
	}
	return fmt.Errorf("no level %d", n)
}

func (g *Game) levelNumber(room string) int {
	var n int
	if g.Order != nil {
		if _, err := fmt.Sscanf(room, "pack:%d", &n); err == nil {
			return n
		}
		return 0
	}
	if _, err := fmt.Sscanf(room, "nivelo%d", &n); err == nil {
		return n
	}
	return 0
}

// levelStart is the score and lives a level began with (the pause menu's
// "Rekomenci la nivelon" puts them back, so a restart can't farm points).
type levelStart struct{ n, score, lives int }

// Level is the level being played (0: a title, rules or highscore room).
func (g *Game) Level() int {
	if g.Waiting() || g.World == nil {
		return 0
	}
	return g.levelStart.n
}

// RestartLevel starts the current level again with the score and lives it
// began with (the pause menu; a remake addition, docs/UI.md).
func (g *Game) RestartLevel() error {
	ls := g.levelStart
	if ls.n == 0 {
		return fmt.Errorf("no level")
	}
	g.G.Score, g.G.Lives = ls.score, ls.lives
	return g.EnterLevel(ls.n)
}

// EnterBonusFrom starts level n and goes to the bonus level from it, as the
// diamond does (the level kept; dying in the bonus level returns to it) —
// the title's level pick, once the bonus level was found.
func (g *Game) EnterBonusFrom(n int) error {
	if err := g.EnterLevel(n); err != nil {
		return err
	}
	g.saved = g.World
	return g.enter("bonusnivelo")
}
