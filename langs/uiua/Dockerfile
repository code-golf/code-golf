# FROM rust:1.78-alpine3.19 as builder

# RUN apk add --no-cache libx11-dev 

# RUN curl -L -o uiua.zip https://github.com/uiua-lang/uiua/releases/download/0.10.3/uiua-bin-x86_64-unknown-linux-gnu.zip \
#     && unzip -j uiua.zip -d /usr/bin/uiua \
#     && rm uiua.zip

FROM rust:1.78-bookworm as builder

RUN apt-get update && apt-get install -y libx11-dev 

RUN cargo install uiua

# breaks with cg base atm
FROM debian:bookworm

COPY --from=0 /usr/local/cargo/bin/uiua /usr/local/bin/uiua

ENTRYPOINT ["uiua"]

CMD ["--version"]
