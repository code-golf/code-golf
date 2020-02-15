FROM rust:1.41.0-alpine3.11 as builder

COPY /compile_and_run.sh /usr/bin/rust

ENTRYPOINT [ "rustc", "--version" ]
