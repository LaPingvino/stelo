# Stelo: recovered game source (from stelo.exe)

Recovered on 2026-11-01 from `stelo.exe` (md5
02849881484c40680e428adf2cc20d4d; the 2018 itch upload of a game whose project
was last saved on 2006-08-22). The engine is the **Game Maker 5.3 runner** with
data format **530** (5.3/5.3a).

## How it was extracted

```
python3 -I gm53_reader.py stelo.exe .   # -> project.json + media
python3 -I gm53_export.py .             # -> objects/ timelines/ rooms/ contact sheet
```

Both scripts are in [gm53-reader](https://github.com/LaPingvino/gm53-reader) (Python 3 + Pillow).
Format sources that were read: LateralGM `GmFileReader.java` (the .gmd 5.3a
layout) and a GM 4.x decompiler gist (the cipher).

- At offset 1,500,000 the exe has `u32 1230500`, then `u32 seed`.
- Every following byte is a substitution through the inverse of
  `table = identity; for i in 1..10000: j = 1 + (i*seed) % 254; swap(t[j], t[j+1])`.
- The decrypted stream holds `1234321, 530` and then a body laid out like a
  5.3a `.gmd`: settings, sounds, sprites, backgrounds, paths, scripts, data
  files, timelines, objects, rooms, game info, room order.

## What's here

| Path | Contents |
|---|---|
| `project.json` | everything parsed, raw (all fields, ids) |
| `objects/*.gml` | 20 objects, every event as readable pseudo-GML. Drag-and-drop actions appear as `action_*()` calls; code actions are verbatim |
| `timelines/*.gml` | 2 timelines (monster spawners) |
| `rooms/*.json` | 9 rooms: instances with object names and positions, views, creation code |
| `rooms/*.png` | each room rendered 1:1 with the original sprites. Invisible helper objects (`marker` = red square, `life_controller`, spawners) are drawn too |
| `rooms/*.txt` | ASCII grid on 32 px cells, with a legend |
| `sprites/`, `backgrounds/` | all 19 sprites and 3 backgrounds as PNG (bottom-left pixel = transparent colour, as in GM) |
| `sounds/` | 4 WAV effects and 5 MIDI tunes |
| `game.ico`, `game_information.rtf` | icon; the F1 game info is empty |
| `sprites_contact.png` | contact sheet of the original sprites |

## Known gaps

- The last ~6 KB after the room order are not parsed. This is probably the
  resource-tree folder layout, which carries no gameplay.
- Action names come from the action library's function names stored in the
  exe. Structural actions (block, else, exit), comments, variables and code
  have no function name and are rendered by kind.
- The `KeyPress <space>` event of `character` creates `<object 32>`, an object
  that no longer exists. It is dead code because `ammo` is always 0.
- No room editor/tree folder names are recovered, and no `.gmd` was rebuilt.
