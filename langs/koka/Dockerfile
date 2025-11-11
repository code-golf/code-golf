ARG VER=3.2.2

FROM node:25.1.0-trixie-slim AS builder

ARG VER

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc

RUN curl -#L https://github.com/koka-lang/koka/releases/download/v$VER/koka-v$VER-linux-x64.tar.gz \
  | tar xz --directory /usr/local

COPY koka.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/koka -s koka.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/x86_64-linux-gnu/libatomic.so.1  \
              /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libgmp.so.10    \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/libstdc++.so.6  /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/gcc /usr/bin/koka            /usr/bin/
COPY --from=0 /usr/local/bin/koka                   \
              /usr/local/bin/node                   /usr/local/bin/
COPY --from=0 /usr/local/share/koka/v$VER/lib/std   /usr/local/share/koka/v$VER/lib/std

ENTRYPOINT ["koka"]

CMD ["--version"]
