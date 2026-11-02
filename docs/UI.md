# A modern UI for Stelo and Pacmanio — APPROVED (Joop 2026-11-02)

Joop (2026-11-02): "fully modernize the UI of both: replace the lettering etc with modern proper fonts, less
blurring etc, remove my old cringe texts and put a nicer explanation, nicer and proper menus similar to Pits
etc etc" — "I used some trickery to get things work at the time, we can do proper menus, level picking etc now."

The levels, rules and timing stay the 2006 game's (the rules docs remain the truth for play). Everything
around play is redesigned; both games get the same UI, each in its own colours.

## 1. Sharp, not blurred
- Today both games draw a fixed 640×480 picture that Ebitengine stretches to the window with smoothing — the
  blur. Instead: the playfield is scaled by whole steps with no smoothing (pixel art stays crisp; a small black
  border where the window doesn't divide evenly; on a phone a fractional step, still unsmoothed), and **all
  text and menus are drawn at the screen's real resolution**, so letters are sharp at any size.

## 2. Fonts
- **Noto Sans** (SIL Open Font License; on this machine; all of ĉ ĝ ĥ ĵ ŝ ŭ) for menus, explanations and the HUD,
  bold for headings and numbers. Replaces the Windows-era fonts (MS Sans Serif, Arial, Courier, Trebuchet) and
  their stand-ins.
- Optional: a rounder display font for the game titles only (e.g. Nunito or Fredoka, also OFL) — a download
  from Google Fonts; only if Joop wants it.

## 3. Menus (as in Pits: arrows + Enter, mouse, and tap a line)
**Title screen** — the game's name drawn in the new font over the game's art, and a menu:
- **Ludi** (play from level 1) · **Daŭrigi** (continue at the furthest level reached; the resume Joop wanted) ·
  **Elekti nivelon** (level picking: a grid of the levels reached so far, each with a small picture of the
  level; in Stelo also the bonus level once found) · **Kiel ludi** (the new explanation) · **Plej bonaj**
  (highscores) · **Agordoj** (language Esperanto / English; sound on / off, music on / off; Stelo: new / 2006
  art) · **Pri** (credits) · desktop only: **Fini**.
- **Pause menu** (Esc, P, or a pause button on phones): Daŭrigi · Rekomenci la nivelon · Agordoj · Ĉefmenuo.
- Stelo's level packs and editor stay reachable from the title menu (Niveloj de aliaj / Redaktilo).

## 4. Texts
- The 2006 rules room (Stelo's klarigo) and thanks room (danklisto), Pacmanio's title subtitle and the
  "Programisto: Joop KIEFTE" corner box go. In their place:
  - **Kiel ludi**: one screen per game — the goal in two sentences, the controls (keys; the phone pads),
    and a legend of the things in the game, each with its own sprite (coin, heart, diamond, checkpoint,
    exit, monsters / ghosts, power star, fruit).
  - **Pri**: short and warm — made by Joop in 2006 with Game Maker, remade in 2026; the new art by @lucasserapiao; music and sound sources; fonts and licences.
- New texts in Esperanto first, English alongside (the translation system Stelo has, extended to Pacmanio).
  Draft texts come to Joop for a read before they go in.

## 5. In play
- **HUD**: one slim bar — score, lives as icons, the level's name — instead of Stelo's cream box and Pacmanio's
  four corner boxes.
- **Level start / end**: the level's name shown briefly at its start ("Nivelo 2 — …"), a short "Nivelo
  finita!" at its end (within the existing pause, so timing doesn't change).
- **Highscores**: a proper table in the new font; name entry with a text field (and the phone keyboard).

## 6. How it gets built (after the yes)
1. Shared UI package (menu widget, text drawing, layout at real resolution, the sharp scaler) — written once,
   used by both games (copied, as the gm engine core is today; later a GameViber library).
2. Stelo first (it has packs, editor and translation already), then Pacmanio.
3. Each step checked by headless screenshots at desktop and phone sizes; the bot playthroughs keep proving
   the rules are untouched.

## Joop's answers (2026-11-02)
1. Title: "make it new" — a new title screen from the game's own sprites and the font, not the 2006 picture.
2. Title font: "whatever looks nice" — Noto Sans Black for titles, Noto Sans for the rest (both installed, OFL;
   no download).
3. Levels: "resume wherever you got, like with Pits" — Daŭrigi continues at the furthest level reached; the
   level picker offers only levels reached.
4. The new art (Stelo) is by @lucasserapiao.

## Questions for Joop (answered above)
1. Title pictures: keep the 2006 title art (Stelo's stelo_komenco, Pacmanio's brick logo) behind the new menu,
   or a new title made from the game's own sprites and the font?
2. A rounder display font for the titles (download), or Noto Sans throughout?
3. Level picking: only levels already reached (as above), or every level open from the start?
4. (ANSWERED 2026-11-02: @lucasserapiao) The commissioned art's artist: name for the credits?
