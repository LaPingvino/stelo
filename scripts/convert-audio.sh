#!/usr/bin/env bash
# convert-audio.sh — the 2006 sounds into formats for today's machines, in audio/effects and audio/music:
#   effects (*.wav)  → Ogg Opus at 48 kHz, Opus' own rate (no resampling at run time), 96 kbit/s
#   music (*.mid)    → rendered with FluidSynth and the FluidR3 GM soundfont (MIT), then Ogg Opus at
#                      64 kbit/s (the web download stays small),
#                      the silence at the end trimmed so loops join
# Played with gopus (github.com/LaPingvino/gopus), a pure-Go Opus decoder. The originals stay as
# they are; this only writes audio/. Needs ffmpeg (with libopus) and fluidsynth.
#
#   scripts/convert-audio.sh [SRC_DIR] [OUT_DIR]
set -euo pipefail
SRC=${1:-original/sounds}
OUT=${2:-audio}
SF2=${SF2:-/usr/share/soundfonts/FluidR3_GM.sf2}
mkdir -p "$OUT/effects" "$OUT/music"
tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT
for f in "$SRC"/*.wav; do
	n=$(basename "$f" .wav)
	ffmpeg -hide_banner -loglevel error -y -i "$f" -ar 48000 -c:a libopus -b:a 96k "$OUT/effects/$n.opus"
done
for f in "$SRC"/*.mid; do
	n=$(basename "$f" .mid)
	fluidsynth -ni -q -g 0.8 -r 48000 -F "$tmp/$n.wav" "$SF2" "$f" >/dev/null
	ffmpeg -hide_banner -loglevel error -y -i "$tmp/$n.wav" \
		-af "areverse,silenceremove=start_periods=1:start_threshold=-60dB,areverse" \
		-ar 48000 -c:a libopus -b:a 64k "$OUT/music/$n.opus"
done
ls -la "$OUT"/effects "$OUT"/music
