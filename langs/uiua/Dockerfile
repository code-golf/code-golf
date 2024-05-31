FROM rust:1.78-bookworm as builder

RUN apt-get update && apt-get install -y libx11-dev 

RUN cargo install uiua

# TODO: breaks if this uses lang-base instead
FROM debian:bookworm

COPY --from=0 /usr/local/cargo/bin/uiua /usr/local/bin/uiua

ENTRYPOINT ["uiua"]

CMD ["--version"]
