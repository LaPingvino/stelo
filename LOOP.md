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
2. [x] docs/LEVEL-FORMAT.md + `cmd/convert`: original/rooms/*.json → levels/*.txt, with a round-trip test.
3. [x] Go module, Ebitengine skeleton (640×480 logical, 30 TPS), loads a level and draws it with the original
   sprites; GitHub Actions → Pages (wasm).
4. [~] GM engine core (internal/gm), needed whether the objects are hand-ported or run as GML: done —
   motion+gravity, inclusive bounding boxes, place_free/place_meeting, move_contact_solid; alarms, the step
   order (alarms, keys, key press, step, move, outside room, collisions with solid put-back, end step),
   instance create/destroy. Next: parent objects (events inherited), room changes and persistence, sleep.
5. [ ] Player: walk, ladders, jump, gravity, solid blocks — level 1 playable.
6. [ ] Coins, hearts, fire/spikes, checkpoints, exit, monsters (walker/flyer, stomp), lives, score, HUD.
7. [ ] All levels, the diamond → bonus level and back, nivelo5 waves.
8. [ ] Title, rules screen, credits, highscores (local storage on web), cheat codes.
9. [ ] Sounds: the 2006 WAVs; MIDI pre-rendered to OGG.
10. [ ] Level editor (in-game mode + its own page) and packs.
11. [ ] New art skin (switchable), roles chosen by trying them.

## Questions for Joop

- 2026-11-01: run the original GML (an embedded GML interpreter + the 76 built-ins the two games use) instead of
  hand-porting the objects? Asked in chat; the engine core (item 4) is needed either way.

- Audio: 2006 sounds/MIDI in the baseline (proposed) or new ones?
- Language: Esperanto only, or also English?

## Blocked

- 2026-11-01: pushing `.github/workflows/pages.yml` is refused — the gh token lacks the `workflow` scope. The
  workflow waits in docs/pages.yml.pending. Joop: `gh auth refresh -h github.com -s workflow`, then
  `git mv docs/pages.yml.pending .github/workflows/pages.yml` and push (Pages is already switched on, source =
  workflow: https://lapingvino.github.io/stelo/).

## Tick log
- 2026-11-01 tick 1 (item 1) — docs/STELO-RULES.md written from original/objects (all 20) + timelines + sprite boxes: player keys/step/collisions, monsters, spawners, lives/HUD, rooms, cheat codes; engine behaviour (GM5 step order, solid put-back) marked [engine]. Found: the ladder grid-snap and Space/ammo are dead code; falling off the bottom costs no life; b_monster ignores invulnerability. NEXT: item 2 (level format + converter).
- 2026-11-01 tick 2 (item 2) — Go module github.com/LaPingvino/stelo; internal/level (Format/Parse, LoadOriginal from original/project.json), cmd/convert → levels/*.txt (9 rooms), docs/LEVEL-FORMAT.md. Grid + separate ladder layer; instance/text lines for off-grid (nivelo1 diamond at 305,308), overlaps, code-carrying and control objects. go test: all 9 rooms round-trip (instances as a set) and levels/ is current. NEXT: item 3 (Ebitengine skeleton + Pages).
- 2026-11-01 tick 3 (item 3) — Ebitengine skeleton: stelo.Files embeds original/ sprites, backgrounds, project.json and levels/; internal/game loads them (GM transparency = bottom-left pixel colour; bbox_lrbt), draws a room in pure Go (colour, tiled background at the view, instances by depth high→low, invisible ones skipped) and has GM view-follow (border + speed, now in the level format and all levels regenerated); cmd/stelo (30 TPS, 640×480, PageUp/Down step rooms until the rules are in); cmd/render → PNG. Test: level 1's first frame = gm53-reader's 1:1 render pixel for pixel (except invisible/sprite-less cells). wasm checked in headless Chromium. Pages workflow (go test, wasm, deploy). NEXT: item 4 (player).
- 2026-11-01 tick 4 (item 4) — internal/gm: Game Maker engine core with tests (inclusive bboxes, place_free, place_meeting, move_contact_solid, gravity+speed motion, direction). ASSUMPTION marked: bbox at the rounded position. Queue renumbered (engine core inserted as 4).
- 2026-11-01 tick 5 (item 4) — internal/gm Step: GM step order, alarms (count down, fire at 0), held/pressed keys, outside-room (box entirely outside), collision dispatch with solid put-back → event → move on if free (ASSUMPTION, marked), end step, destroy at step end; deterministic event order (key code; collisions by name, ASSUMPTION vs GM's object index). Tests: monster_dead's 10-step alarm, landing on a block via the character's own block event, outside-room, keys.
