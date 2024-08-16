FROM rust:1.80.1-slim-bookworm AS builder

RUN apt-get update \
 && apt-get install -y curl

RUN curl -L https://github.com/betaveros/noulith/archive/828de52.tar.gz | tar xz --strip-components=1

RUN cargo build --release \
 && cargo fix --allow-no-vcs --lib --package noulith

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/ld-linux-x86-64.so.2 \
              /lib/x86_64-linux-gnu/libc.so.6            \
              /lib/x86_64-linux-gnu/libgcc_s.so.1        \
              /lib/x86_64-linux-gnu/libm.so.6            /lib/x86_64-linux-gnu/
COPY --from=0 /lib64                                     /lib64
COPY --from=0 /target/release/noulith                    /usr/bin/

ENTRYPOINT ["noulith"]
