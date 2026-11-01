# Stelo remake — loop

Work one bounded chunk (~10 min) per tick, top unchecked item first. Design: [docs/DESIGN.md](docs/DESIGN.md)
(approved by Joop 2026-11-01). The 2006 game is in [original/](original/) (recovered source: objects/*.gml,
rooms/*.json, sprites, sounds).

## Non-negotiables

- The original is the truth: rules, numbers and level layouts come from `original/` (objects' code, rooms),
  not from guessing. Where the code is unclear, say so in the rules doc and ask Joop (no Crostini needed;
  wine is installed locally if a check is really needed, ~/Projects/stelo/win holds the exe).
- `win/`, `linux/` and any `.exe` never go into git (they are the itch builds).
- New art: `~/Projects/stelo-remake/research/assets/` (Joop's commissioned 32×32 set). Which sprite plays which
  role is ours to choose (Joop: the artist was unsure); record the choices in docs/DESIGN.md.
- Commit and push to origin main each tick; the web build goes to GitHub Pages.
- Go + Ebitengine; logic at the original 30 steps per second.

## Queue

1. [x] docs/STELO-RULES.md from original/objects/*.gml and timelines: every object, every event, in step order
   (movement numbers, gravity, ladders, stomp test, scores, lives, checkpoints, bonus level, spawners, cheat codes).
2. [ ] docs/LEVEL-FORMAT.md + `cmd/convert`: original/rooms/*.json → levels/*.txt, with a round-trip test.
3. [ ] Go module, Ebitengine skeleton (640×480 logical, 30 TPS), loads a level and draws it with the original
   sprites; GitHub Actions → Pages (wasm).
4. [ ] Player: walk, ladders, jump, gravity, solid blocks — level 1 playable.
5. [ ] Coins, hearts, fire/spikes, checkpoints, exit, monsters (walker/flyer, stomp), lives, score, HUD.
6. [ ] All levels, the diamond → bonus level and back, nivelo5 waves.
7. [ ] Title, rules screen, credits, highscores (local storage on web), cheat codes.
8. [ ] Sounds: the 2006 WAVs; MIDI pre-rendered to OGG.
9. [ ] Level editor (in-game mode + its own page) and packs.
10. [ ] New art skin (switchable), roles chosen by trying them.

## Questions for Joop

- Audio: 2006 sounds/MIDI in the baseline (proposed) or new ones?
- Language: Esperanto only, or also English?

## Blocked

## Tick log
- 2026-11-01 tick 1 (item 1) — docs/STELO-RULES.md written from original/objects (all 20) + timelines + sprite boxes: player keys/step/collisions, monsters, spawners, lives/HUD, rooms, cheat codes; engine behaviour (GM5 step order, solid put-back) marked [engine]. Found: the ladder grid-snap and Space/ammo are dead code; falling off the bottom costs no life; b_monster ignores invulnerability. NEXT: item 2 (level format + converter).
