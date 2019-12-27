FROM rust:1.40.0-alpine3.10 as builder

COPY /compile_and_run.sh /usr/bin/rust

ENTRYPOINT [ "rustc", "--version" ]
