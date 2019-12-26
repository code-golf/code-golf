#!/bin/sh -e
src="$1"
shift
export RUST_BACKTRACE=1
export RUSTUP_HOME=/usr/local/rustup
export PATH=/usr/local/cargo/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
/usr/local/cargo/bin/rustc -o /tmp/rs-code "$src"
exec -a "$src" /tmp/rs-code "$@"
