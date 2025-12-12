ARG VER=9.6 PATCH=3

FROM debian:trixie-slim AS builder

ARG VER PATCH

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl

WORKDIR /j$VER

RUN curl -#L https://www.jsoftware.com/download/j$VER/install/j$VER.${PATCH}_linux64.tar.gz \
  | tar xz --exclude updateje.sh --strip-components 1

RUN ./bin/install-usr.sh

FROM codegolf/lang-base

ARG VER

COPY --from=0 /etc/j/$VER/profile.ijs               /etc/j/$VER/
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libgmp.so.10    \
              /lib/x86_64-linux-gnu/libj.so.$VER    \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/ijconsole                    /usr/bin/j
COPY --from=0 /usr/share/j/$VER/addons/dev          /usr/share/j/$VER/addons/dev
COPY --from=0 /usr/share/j/$VER/system              /usr/share/j/$VER/system

ENV GLIBC_TUNABLES='glibc.rtld.execstack=2'

ENTRYPOINT ["j"]

CMD ["-js", "exit echo JVERSION"]
