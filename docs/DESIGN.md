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
2. The game playing level 1, then all levels and the bonus level; checked side by side with the original running under Wine.
3. Title, rules, credits, highscores; sounds (2006 WAVs, MIDI pre-rendered).
4. The editor.
5. The new art skin.

## Decisions and open questions

- Audio: the 2006 sounds and MIDI, or new ones? (Proposed: the 2006 ones in the baseline.)
- Sprite roles: Joop leaves them to us (the artist was unsure of the roles); try the new sprites in each role and pick what reads best.
- Language: Esperanto only, or also English?
- Cheat codes: KEEP (Joop, 2026-11-01). Credits screen: proposed keep.
- DONE: recovered source published at github.com/LaPingvino/stelo (original/), reader at github.com/LaPingvino/gm53-reader (MIT).
