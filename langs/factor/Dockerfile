FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc

ENV VER=0.100

RUN curl -#L https://downloads.factorcode.org/releases/$VER/factor-linux-x86-64-$VER.tar.gz \
  | tar xz --directory /usr/local/bin --strip-components 1

COPY factor.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/factor -s factor.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/librt.so.1      \
              /lib/x86_64-linux-gnu/libstdc++.so.6  /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/factor                       /usr/bin/
COPY --from=0 /usr/local/bin/factor                 \
              /usr/local/bin/factor.image           /usr/local/bin/
COPY --from=0 /usr/local/bin/basis                  /usr/local/bin/basis
COPY --from=0 /usr/local/bin/core                   /usr/local/bin/core

ENTRYPOINT ["factor"]

CMD ["--version"]
