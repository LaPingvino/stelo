# How Stelo works

Everything below is read from the recovered 2006 project in [`original/`](../original/): the objects' code in
`original/objects/*.gml`, the timelines, the rooms and the sprite data in `original/project.json`. The remake is
written from this document. Where Stelo relies on how **Game Maker 5.3 itself** behaves (the order of events in a
step, what a "solid" does), that is marked **[engine]**: it comes from Game Maker's documented behaviour, not
from Stelo's code.

## 1. Time and space

- Every room runs at **30 steps per second** (room speed 30 in all nine rooms).
- Coordinates are pixels; levels are 640×480 (one screen) or 2560×480 (levels 2–4, the view follows the player).
- Every sprite is 32×32 with its origin at the top-left corner, except `sprite_life` (10×10). Collisions use each
  sprite's **bounding box** (not its pixels), given here as left–right, top–bottom inside the 32×32 cell:

| Sprite | Used by | Bounding box (x, y) |
|---|---|---|
| `sprite_block` | block | 0–31, 0–31 |
| `sprite_right`, `sprite_left`, `sprite_climbing` | the player | 7–23, 2–31 |
| `sprite_monsterr`, `sprite_monsterl` | monster, b_monster | 3–28, 1–31 |
| `sprite_monsterflat` | monster_dead | 2–29, 27–31 |
| `sprite_flyingr`, `sprite_flyingl` | flyer | 4–28, 8–24 |
| `sprite_ladder` | ladder | 6–25, 0–31 |
| `sprite_death` | death (fire) | 2–30, 4–31 |
| `sprite_stelo_poentoj` | star coin | 0–31, 0–31 |
| `sprite_bonus` | heart | 3–28, 4–26 |
| `transirejo_na`, `transirejo_a` | checkpoint (unused / used) | 5–24, 2–28 |
| `diamanto` | diamond | 1–29, 3–29 |
| `sprite_level` | level exit (and the invisible spawners) | 1–31, 0–31 |
| `sprite_marker` | marker (invisible) | 0–31, 0–31 |

- **Solid** objects: only `block`.

## 2. One step [engine]

Game Maker 5 runs each step in this order, for every instance:

1. alarms;
2. keyboard events (a key held down), then key-press events;
3. the step event;
4. movement: `vspeed += gravity` (in the gravity direction), then `x += hspeed`, `y += vspeed`;
5. "outside room" events;
6. collision events. When the other instance is **solid**, the moving instance is first put back where it was
   before step 4, then the collision event runs;
7. end-step events;
8. drawing.

Terms used below: **free(dx, dy)** = the player's box moved by (dx, dy) overlaps no solid (`place_free`);
**meets(obj)** = the box at the current place overlaps an instance of obj; **move to contact(dir, max)** =
move pixel by pixel in that direction, up to max pixels, stopping before a solid.

## 3. The player (`character`)

Starts with `netusxu = 0` (invulnerability off).

**Keys, while held:**

- **Right**: face right. If free(+4, 0), move 4 right; else if free(+4, −8), move (+4, −8) then move to contact
  downwards, max 8 (this walks up a step of up to 8 px). **Left** is the same mirrored.
- **Up**: on a ladder (meets ladder): if free(0, −3), move 3 up. Not on a ladder: if standing (not free(0, +1)),
  jump: `vspeed = −10`.
- **Down**: on a ladder: if free(0, +3), move 3 down.
- **Space** (key press): would create object 32 if `ammo > 0`; `ammo` is always 0 and object 32 no longer exists,
  so this does nothing.

**Each step:**

1. The room background scrolls: `background_hspeed = −1`, or `+1` when the player faces left.
2. Gravity: if free(0, +1) (in the air) gravity is 0.5 downwards, else 0.
3. `vspeed` is capped at 12 (falling).
4. On a ladder (meets ladder): gravity 0, `vspeed = 0`, sprite becomes `sprite_climbing` (it stays so until Left
   or Right sets the walking sprite again).
5. A snap to the 32 px grid when standing still while the sprite is `sprite_ladder`: the player never has that
   sprite, so this never happens.
6. Returning from the bonus level: if the room is persistent (it was left for the bonus level), make it
   non-persistent again, stop all sounds and loop `muziko1`.
7. Cheat codes: see §8.

**Collisions** (all except block and the diamond first check: if invulnerable, nothing happens):

| With | What happens |
|---|---|
| `block` (solid) | [engine] put back first; then move to contact in the direction of motion, max 12; `vspeed = 0` |
| `stelo_poentoj` (star) | sound `getmushroom`, score +10, the star disappears |
| `bonus` (heart) | sound `finishlevel`, lives +1, the heart disappears |
| `transirejo` (checkpoint) | if not used yet: the player's start point (`xstart, ystart`) becomes the checkpoint's position, it is marked used and shows `transirejo_a` |
| `death` (fire) | sound `killcharacter`, the game pauses 1 s, lives −1, back to the start point |
| `monster` | if falling (`vspeed > 0`) and `y < monster.y + 8`: stomp — sound `killmonster`, score +50, the monster stops and becomes `monster_dead`. Otherwise: sound `killcharacter`, pause 1 s, lives −1, back to the start point |
| `flyer` | sound `killcharacter`, pause 1 s; in the bonus level: back to the saved level (no life lost); elsewhere: lives −1, back to the start point. Flyers can't be stomped |
| `b_monster` (bonus level) | stomp as for `monster` but score +10; otherwise sound `killcharacter`, pause 1 s, back to the saved level (no life lost) |
| `iru_al_bonuso` (diamond) | the current room becomes persistent (kept as it is), `global.ludnivelo` = this room, the diamond disappears, sounds stop, `bonussono` loops, go to `bonusnivelo` |
| `levelexit` | sound `finishlevel`, pause 1 s, go to the next room (if there is one) |

**Outside the room:** when `y > room_height` (fell off the bottom): sound `killcharacter`, pause 1 s, back to the
start point. **No life is lost.**

The start point is where the player was placed in the room, until a checkpoint changes it.

## 4. Monsters

- **`monster`** (walker): starts moving right at speed 4 (`hspeed = 4`). No gravity. Hitting a `marker`: back to
  its previous position and reverse; hitting a block: reverse [engine: put back first]. End of step: sprite
  `sprite_monsterr` when `hspeed > 0`, else `sprite_monsterl`.
- **`b_monster`** (bonus level): as `monster`, plus gravity like the player (0.5 when free(0, +1), else 0); on a
  block: move to contact (direction of motion, max 12), `vspeed = 0`, reverse.
- **`flyer`**: starts moving right at speed 3; no gravity; reverses on a marker or a block (without stepping back
  on a marker). Sprite by direction, as above.
- **`monster_dead`**: shows `sprite_monsterflat` and disappears after 10 steps.
- **`marker`**: invisible; only turns monsters around.

## 5. Spawners

- **`monsterfarilo`** (bonus level, invisible): each step, if there is no `b_monster` left and its timeline isn't
  running, it starts timeline `faruMonsterojn`: at step 30, 60 and 90 a `b_monster` appears at the spawner's
  position; after the third the spawner can start again.
- **`finmonsterfarilo`** (level 5, invisible): the same with `monster` and timeline `faruFinmonsterojn`.

## 6. Lives, score and the HUD (`life_controller`)

- A persistent object (it lives across rooms). At its first room start (`global.startita` false) it sets lives
  to 3; `global.startita` is set to false once, at game start.
- **No more lives** → room `finlisto` (the highscore screen).
- **Draws** each step: a filled rectangle from (view left + 5, view top + 5) to (view left + 100, view top + 25)
  in colour RGB (231, 230, 180); the score as "Poentoj: N" at (view left + 10, view top + 5); and one
  `sprite_life` per life in a row starting at (player x, player y − 10).
- Score starts at 0.

## 7. Rooms and screens

Room order: `klarigo` → `nivelo1` → `nivelo2` → `nivelo3` → `nivelo4` → `nivelo5` → `danklisto` → `finlisto`;
`bonusnivelo` is reached only through the diamond (or a cheat).

- **Game start** (`klarigilo`): plays `komenco`, draws the title background `stelo_komenco`, waits for a key, then
  stops `komenco` and loops `muziko1`.
- **`klarigo`** (rules screen): drawn, waits for a key, next room. Its text is in `tekstilo` instances (each draws
  its `teksto` in green, Arial Unicode MS 11, 3 px right of its position).
- **`nivelo1`–`nivelo5`**: the levels. Level 5 has the `finmonsterfarilo` spawners.
- **`bonusnivelo`**: endless; dying there returns to the saved level, which continues exactly as it was left (it
  was made persistent). Its exit is death: there is no `levelexit`.
- **`danklisto`**: thanks/credits; reached from level 5's exit.
- **`finlisto`** (`plejbonaj_listo`): stops the life controller, loops `finavenko`. If the score beats the 10th
  highscore, asks for a name ("Vi tiom bone ludis ke vi aperos en la listo de la plej bonaj ludantoj! Via nomo:")
  until Enter and adds it; otherwise shows "Ho ve... Vi ne sufiĉe bone ludis por aperi en la listo de la plej
  bonaj ludantoj..." and waits for a key. Then draws the highscore table at (x + 5, y + 5)–(x + 480, y + 320),
  waits for a key and restarts the game.

## 8. Cheat codes

Typed at any time during a level (checked each step against everything typed; a match clears the typed text):

| Code | Effect |
|---|---|
| `sekvosio` | next room |
| `vivuplilonge` | lives +1 |
| `netusxumin` | invulnerable (all the checks above) |
| `mibedauxras` | invulnerable off; music back to `muziko1` |
| `esperanto` | music: `LaEspero` (La Espero) |
| `fekalfajro` | every fire becomes a block |
| `miamasGrandanFraton` | score +500 |
| `pasispasinteco` | clears the highscores; score −1 |
| `pliajpoentoj` | goes to the bonus level, as the diamond does |

## 9. Open points

- The 1 s pause (`action_sleep(1000)`) freezes the whole game in the original, drawing included. The remake
  keeps its length and the logic waits exactly as before, but the screen stays alive (quality of life, Joop
  2026-11-02: "keep the beat, no freeze"): the moment the pause began is shown with its animations running,
  the view and the HUD hold still; the rest of the event shows when the pause ends.
- Highscore table drawing details (Game Maker's `draw_highscore` layout and font) are engine behaviour; to be
  matched by eye against the original.
