#!/bin/sh -e

ttf='build/Twemoji Mozilla.ttf'
woff='build/Twemoji Mozilla.woff2'
final='build/twemoji.woff2'

# Make the font multiple times as the compressed size in non-deterministic.
for i in $(seq 1 100); do
    rm -f "$ttf"
    make
    woff2_compress "$ttf"

    # If $final doesn't exist, or $woff is smaller, copy $woff to $final.
    if [ ! -f "$final" ] || { [ `stat -c %s "$woff"` -lt `stat -c %s "$final"` ]; }; then
        cp "$woff" "$final"
    fi
done
