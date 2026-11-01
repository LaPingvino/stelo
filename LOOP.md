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
4. [x] GM engine core (internal/gm), needed whether the objects are hand-ported or run as GML: done —
   motion+gravity, inclusive bounding boxes, place_free/place_meeting, move_contact_solid; alarms, the step
   order (alarms, keys, key press, step, move, outside room, collisions with solid put-back, end step),
   instance create/destroy; parent objects (handlers inherited, collisions and with() cover children). room Goto request, action_sleep (frozen steps), image_index animation, instance variables; the game now runs levels on gm (instances created with create events, drawn from gm). Next: room changes/persistence in the game loop, then item 5.
5. [~] Player (hand-ported from character.gml; Joop hasn't answered the GML question yet): keys, step, block, outside room done. walk, ladders, jump, gravity, solid blocks — level 1 playable.
6. [x] Coins, hearts, fire/spikes, checkpoints, exit, monsters (walker/flyer, stomp), lives, score, HUD. Done: coins, hearts, fire, checkpoints, monsters/flyers/monster_dead, stomp, lives/score (gm.Globals). Room changes done (Game: @next, named, @bonus keeps the level, @saved returns; no more lives → finlisto). HUD done. Next: item 7.
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
- 2026-11-01 tick 6 (item 4) — parent objects in internal/gm: IsA, handler lookup up the chain (create/alarm/step/outside/end step; key and collision maps merged, child wins), collisions and All()/with() include children. Test with Pacmanio's monster→scared→intermonster chain; it records that both of pacman's monster/scared events see an intermonster — the open [engine] point for Pacmanio §3 stays open.
- 2026-11-01 tick 7 (item 4) — internal/game now runs on internal/gm: levels become gm rooms (objects with sprite box, solidity, depth, visibility), instances created in level order, drawn from gm; World.Step = room step + view. gm: Room.Goto, Sleep/SleepMS (test: 1000 ms = 30 frozen steps), image_index/image_speed, per-instance Vars. Level-1 pixel test still passes; wasm builds. Item 4 done; NEXT item 5 (player) — waits on Joop's GML answer only in HOW the objects are written.
- 2026-11-01 tick 8 (item 5) — the player, hand-ported from character.gml onto internal/gm (internal/game/objects.go): Left/Right walk 4 with the 8 px step-up, Up climbs 3 on a ladder or jumps (vspeed −10) when standing, Down climbs, step gravity/cap 12/ladder (climbing sprite), block collision (move to contact 12, vspeed 0), outside room (1 s, back to start, no life lost). Tests on level 1: stands still, jump ≈105 px and lands at 416, walks 4/step, the start-row ladder at cell 2 catches the box at x 52 (climbing sprite, Up climbs 3/step, hangs) — the original's behaviour, not a bug. Visual check: mid-jump and on the ladder. PARTIAL: background drift, sounds, cheat codes. NEXT: item 6.
- 2026-11-01 tick 9 (item 6) — character's collisions (stelo_poentoj +10, bonus +1 life, transirejo start point + used sprite, death/monster/flyer: sound, 1 s, −1 life, back to start; stomp: falling and y < monster.y+8 → +50, monster stops and becomes monster_dead), all behind netusxu; monster (right at 4, marker: step back + reverse, block: reverse, sprite by hspeed in end step), flyer (3; reverses on marker/block), monster_dead (alarm 10). gm: Globals (score, lives, NoMoreLives, sound hook), Room.Change (action_change_object). Tests: coin, stomp + flattened monster gone after 10 steps, side hit costs a life with the 1 s pause.
- 2026-11-01 tick 10 (item 6) — internal/game.Game runs rooms in RoomOrder with shared gm.Globals; room changes wait for the pause; the diamond keeps the level (persistent) for the bonus level and a flyer/b_monster there returns to it unchanged; levelexit → next room; no more lives → finlisto. FIXED: init-order bug (a wrapper registered before the character's own init was lost) — the character's extra collisions are now wired explicitly. cmd/stelo plays through Game. Test: diamond → bonus → flyer → back to the same nivelo1 (diamond gone, no life lost) → exit → nivelo2.
- 2026-11-01 tick 11 (item 6) — the life controller's HUD: action_color(brush 11855591, pen 0) → light-yellow box (231,230,180) with a black outline at view (5,5)–(100,25); 'Poentoj: N' at (10,5) in black; one sprite_life per life from (player.x, player.y−10). Font: Liberation Sans (OFL, fonts/ with its licence) standing in for Arial Unicode MS 10 pt at 96 dpi. Drawn in every room with a player (= where the persistent life_controller exists). Checked against the Wine screenshot of nivelo1: layout and colours match; text 57×12 px vs the original's 55×10 — Wine likely substituted Arial Unicode MS, so not conclusive (open). Item 6 done; NEXT item 7 (all levels: nivelo5 waves, b_monster, spawners).
