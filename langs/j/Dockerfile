ARG VER=9.6

FROM debian:bookworm-slim AS builder

ARG VER

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl

RUN curl -# https://www.jsoftware.com/download/j$VER/install/j$VER.3_linux64.tar.gz \
  | tar xz

# FIXME Don't run updatejs.sh as it tries to update to AVX512 because my build
# laptop supports it. But prod servers do not. Is there an override flag?
# RUN cd j$VER && ./updateje.sh

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /j$VER/addons/dev                     /opt/j/addons/dev
COPY --from=0 /j$VER/bin/jconsole                   \
              /j$VER/bin/libgmp.so                  \
              /j$VER/bin/libj.so                    \
              /j$VER/bin/profile.ijs                /opt/j/bin/
COPY --from=0 /j$VER/system                         /opt/j/system

ENTRYPOINT ["/opt/j/bin/jconsole"]
