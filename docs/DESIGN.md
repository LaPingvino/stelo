# Stelo remake — design (approved by Joop 2026-11-01)

Modelled on [The Puzzle Pits](https://github.com/LaPingvino/The-Puzzle-Pits): the original game first, faithful,
with its rules and level format written down; an editor on top that makes levels for it.

## 1. Baseline: the original game

- Go + Ebitengine, one codebase for desktop and the browser (WebAssembly on GitHub Pages).
- Play the 2006 game as it was: title, rules screen, levels 1–5, the diamond-entered bonus level, credits, highscores.
  The rooms are converted 1:1 from `original/rooms/*.json`; the bonus level keeps the paused level's state, as in the original.
- Logic runs at the original 30 steps per second with the original numbers (walk 4, climb 3, jump −10, gravity 0.5, terminal 12; scores, lives, checkpoints), so it plays the same.
- `docs/STELO-RULES.md`: the rules as read from the recovered object code (the PITS-RULES.md counterpart): what each object does, step by step, in order. The port is written from it, and it is checked against play under Wine.
- Graphics: the original sprites by default, so the baseline looks like 2006. The new 32×32 art is a second skin, switchable, with gaps filled by the old sprites until drawn.

## 2. Level format and packs

- `docs/LEVEL-FORMAT.md`: a plain-text level file (a header and a 32 px grid of one character per cell; spawners and markers as cells), readable and diffable.
  The original rooms are converted into it once, and a test checks that they convert back to the same instance lists.
- Packs as in Pits: the original levels are the first pack; any other is a folder (or a zip) of level files with a `pack.txt` (name, author). Progress is kept per pack.

## 3. Level editor

- Built into the same program (a mode from the title screen, and its own page on GitHub Pages), so it uses the game's own drawing and rules.
- Opens the original levels or any level file. Paint cells from a palette (walls, floor, ladders, spikes, coins, hearts, diamond, checkpoint, exit, monsters, spawners, markers), set the level size, and use undo.
- **Play** tries the level straight away; save and download a level or a pack (zip).

## 4. Order of work

1. Rules doc and level converter (the original rooms into level files, with the round-trip test).
2. The game playing level 1, then all levels and the bonus level; checked against the recovered rooms and code.
3. Title, rules, credits, highscores; sounds (2006 WAVs, MIDI pre-rendered).
4. The editor.
5. The new art skin.

## Decisions and open questions

- Audio: the 2006 sounds and MIDI, or new ones? (Proposed: the 2006 ones in the baseline.)
- Sprite roles: Joop leaves them to us (the artist was unsure of the roles); try the new sprites in each role and pick what reads best.
- Language: Esperanto only, or also English?
- Cheat codes: KEEP (Joop, 2026-11-01). Credits screen: proposed keep.
- DONE: recovered source published at github.com/LaPingvino/stelo (original/), reader at github.com/LaPingvino/gm53-reader (MIT).

## New art (item 11)

Joop's commissioned 32×32 set is in `art/nova/`, a drawing-only skin: collision boxes stay the original's, so the
rules don't change. It is the default (Joop 2026-11-02: "switch to the commission art first"); F3 switches
to the 2006 art and back; `?original` on the web (STELO_ART=original on the desktop) starts with the 2006 art.
It animates at its own, calmer pace (about 3–6 frames a second, each instance offset so rows don't flash
together; the protagonist only walks its cycle while it moves). The roles (the artist was unsure;
ours to choose, Joop 2026-11-01):

| Original | New art |
|---|---|
| player (sprite_right / left / climbing) | the protagonist star (its left sheet; climbing: its first frame) |
| monster, b_monster | the orange villain (mirrored facing left) |
| flyer | the green villain (it can't be stomped, so it looks meaner) |
| star coin | a star-marked coin: the checkpoint's star (its pole left out), twinkling, on a round blue coin — blue stands out on the wooden backgrounds where amber blended in (Joop 2026-11-02) |
| heart, diamond | heart, diamond |
| checkpoint unused / used | the orb on a pole (a switch) / the checkpoint's star, shining (Joop 2026-11-02: the star in both states hid the change) |
| fire | the spikes |
| level exit | the arrow sign |
| ladder | the ladder |
| block | floor with grass where nothing is above it, the wall elsewhere |

Derived from the set (no art was made for them): the flattened monster (the villain squashed into the bottom 8 px)
and the life icon (the protagonist at 10×10). Not covered (original art kept): the title, the backgrounds.
