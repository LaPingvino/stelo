package game

import (
	"container/heap"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/LaPingvino/stelo"
	"github.com/LaPingvino/stelo/internal/gm"
)

// The play-check (LOOP item 7): a bot shows each level can be played
// through to its exit with the original rules — no life lost, no cheat,
// not via the diamond. Best-first search over held-key choices, each held
// for a few steps, from cloned worlds (World.Clone); states already seen
// (the player's 8 px cell and whether rising) are skipped; nearest to the
// exit first.

var botMoves = []held{
	{}, {gm.KeyRight: true}, {gm.KeyLeft: true}, {gm.KeyUp: true}, {gm.KeyDown: true},
	{gm.KeyRight: true, gm.KeyUp: true}, {gm.KeyLeft: true, gm.KeyUp: true},
}

const botHold = 4 // steps per choice

// botCell is the seen-state grid (px).
var botCell = 8

type botNode struct {
	w     *World
	steps int
	dist  float64
}

type botQueue []*botNode

func (q botQueue) Len() int           { return len(q) }
func (q botQueue) Less(i, j int) bool { return q[i].dist < q[j].dist }
func (q botQueue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *botQueue) Push(x any)        { *q = append(*q, x.(*botNode)) }
func (q *botQueue) Pop() any {
	old := *q
	n := old[len(old)-1]
	*q = old[:len(old)-1]
	return n
}

// playThrough searches until the room asks for the next one; it returns the
// steps taken and the states tried, or -1 when the budget runs out.
func playThrough(t *testing.T, w *World, budget int) (steps, tried int) {
	exits := w.Room.All("levelexit")
	if len(exits) == 0 {
		t.Fatal("no levelexit")
	}
	lives := w.Room.G.Lives
	dist := func(w *World) float64 {
		p := w.Room.All("character")
		if len(p) == 0 {
			return math.Inf(1)
		}
		best := math.Inf(1)
		for _, e := range exits {
			best = min(best, math.Abs(p[0].X-e.X)+math.Abs(p[0].Y-e.Y))
		}
		return best
	}
	seen := map[[3]int]bool{}
	q := &botQueue{{w: w, dist: dist(w)}}
	for q.Len() > 0 && tried < budget {
		n := heap.Pop(q).(*botNode)
		for _, mv := range botMoves {
			c := n.w.Clone()
			for range botHold {
				c.Step(mv)
			}
			tried++
			r := c.Room
			if r.Goto == "@next" {
				return n.steps + botHold, tried
			}
			if r.Goto != "" || r.G.Lives < lives || r.G.NoMoreLives {
				continue // the diamond, or a life lost
			}
			p := r.All("character")
			if len(p) == 0 {
				continue
			}
			rising := 0
			if p[0].VSpeed < 0 {
				rising = 1
			}
			k := [3]int{int(p[0].X) / botCell, int(p[0].Y) / botCell, rising}
			if seen[k] {
				continue
			}
			seen[k] = true
			heap.Push(q, &botNode{w: c, steps: n.steps + botHold, dist: dist(c)})
		}
	}
	return -1, tried
}

func TestPlayCheckLevels(t *testing.T) {
	if testing.Short() {
		t.Skip("the play-check bot takes a while")
	}
	a, err := LoadAssets(stelo.Files)
	if err != nil {
		t.Fatal(err)
	}
	// The wide levels take minutes (nivelo2: 69 330 states, 3½ min), so CI
	// runs 1 and 5; STELO_PLAYCHECK=nivelo2,nivelo3 runs others.
	levels := []string{"nivelo1", "nivelo5"}
	if l := os.Getenv("STELO_PLAYCHECK"); l != "" {
		levels = strings.Split(l, ",")
	}
	for _, name := range levels {
		t.Run(name, func(t *testing.T) {
			w, err := LoadLevel(stelo.Files, a, name)
			if err != nil {
				t.Fatal(err)
			}
			w.Room.G.Lives = 3
			steps, tried := playThrough(t, w, 200000)
			if steps < 0 {
				t.Fatalf("no way to the exit found (%d states tried)", tried)
			}
			t.Logf("exit reached after %d steps (%.1f s of play), %d states tried", steps, float64(steps)/30, tried)
		})
	}
}
