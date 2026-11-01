# Stelo level files

A level is a text file (`levels/<name>.txt`). The original nine rooms were converted from
[`original/project.json`](../original/project.json) with `go run ./cmd/convert`, and a test
(`go test ./internal/level`) checks that every one of them reads back with the same settings and the same
instances, and that `levels/` is up to date.

```
stelo-level 1
name nivelo1
size 640 480
speed 30
colour #00FF80
background background_sky
view 0 0 640 480 follow character border 64 64 speed 4 4
grid
#...|......M...|...#
…one row per 32 px, one character per cell…
ladders
..........H.........
…the same size…
instance iru_al_bonuso 305 308
instance life_controller 32 0
end
```

## Header lines

| Line | Meaning |
|---|---|
| `stelo-level 1` | first line, the format's version |
| `name N` | the room's name |
| `size W H` | size in pixels (a multiple of 32 in every original room) |
| `speed N` | steps per second (30 everywhere) |
| `colour #RRGGBB` | background colour |
| `persistent` | the room is kept as it is when left (none of the originals; the bonus level does this at run time) |
| `background NAME` | a background image, tiled across the room |
| `caption "…"` | the window caption, as a Go-quoted string |
| `view X Y W H [follow OBJECT border H V speed H V]` | the visible part of the room, following an object: it scrolls when the object comes within `border` pixels of an edge, by at most `speed` pixels per step (−1: no limit) |
| `code "…"` | the room's creation code (Game Maker code, kept as text; the remake implements what the originals use) |

## Grid and ladders

`grid` is followed by H/32 rows of W/32 characters; `.` is empty:

| | Object | | Object |
|---|---|---|---|
| `#` | block (wall) | `*` | star (`stelo_poentoj`) |
| `@` | the player (`character`) | `+` | heart (`bonus`) |
| `M` | monster | `D` | diamond (`iru_al_bonuso`) |
| `F` | flyer | `T` | checkpoint (`transirejo`) |
| `b` | bonus-level monster (`b_monster`) | `E` | level exit |
| `X` | fire (`death`) | `\|` | marker (turns monsters) |
| `S` | spawner, bonus level (`monsterfarilo`) | `s` | spawner, level 5 (`finmonsterfarilo`) |

`ladders` is a second layer of the same size holding only `H`, so a ladder can share a cell with anything on the
grid (stars on ladders are common).

## Other instances

After the two layers, one line each:

- `instance OBJECT X Y ["code"]` — anything the grid can't hold: off the 32 px grid (level 1's diamond at
  305, 308), a second object in an occupied cell, an object with no grid character (`life_controller`,
  `klarigilo`, `plejbonaj_listo`), or an instance with creation code.
- `text X Y "…"` — a text label (`tekstilo` whose creation code is `teksto='…';`). `#` is a line break, as in
  Game Maker.

`end` closes the file.

## Not kept

The order in which instances were created in the original. It only matters where two instances' events depend
on each other's order within a step, which the remake handles by its own fixed order.
