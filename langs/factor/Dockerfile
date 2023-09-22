FROM debian:bookworm-slim AS builder

RUN apt-get update && apt-get install -y curl

RUN curl https://downloads.factorcode.org/releases/0.99/factor-linux-x86-64-0.99.tar.gz \
  | tar xz

FROM codegolf/lang-base

COPY --from=0 /factor/basis                         /factor/basis
COPY --from=0 /factor/core                          /factor/core
COPY --from=0 /factor/factor                        \
              /factor/factor.image                  /factor/
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/librt.so.1      \
              /lib/x86_64-linux-gnu/libstdc++.so.6  /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/

ENTRYPOINT ["/factor/factor"]

CMD ["--version"]
