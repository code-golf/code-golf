#!/bin/sh -e
src="$1"
shift
/usr/local/cargo/bin/rustc -o /tmp/rs-code "$src"
RUST_BACKTRACE=1 exec -a "$src" /tmp/rs-code "$@"
