FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl

ENV VER=0.100

RUN curl -#L https://downloads.factorcode.org/releases/$VER/factor-linux-x86-64-$VER.tar.gz \
  | tar xz --directory /usr/bin --strip-components 1

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/librt.so.1      \
              /lib/x86_64-linux-gnu/libstdc++.so.6  /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/factor /usr/bin/factor.image /usr/bin/
COPY --from=0 /usr/bin/basis                        /usr/bin/basis
COPY --from=0 /usr/bin/core                         /usr/bin/core

ENTRYPOINT ["factor"]

CMD ["--version"]
