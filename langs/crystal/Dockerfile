FROM crystallang/crystal:1.7.1-alpine as builder

RUN strip /usr/bin/crystal

FROM codegolf/lang-base

ENV CRYSTAL_CACHE_DIR=/tmp PATH=/usr/bin:/bin

COPY --from=0 /bin/sh                              /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1             \
              /lib/libz.so.1                       /lib/
COPY --from=0 /usr/bin/cc /usr/bin/crystal         /usr/bin/
COPY --from=0 /usr/include                         /usr/include
COPY --from=0 /usr/lib                             /usr/lib
COPY --from=0 /usr/libexec/gcc                     /usr/libexec/gcc
COPY --from=0 /usr/local/lib/libgc.a               /usr/local/lib/
COPY --from=0 /usr/share/crystal                   /usr/share/crystal
COPY --from=0 /usr/x86_64-alpine-linux-musl/bin/ld /usr/x86_64-alpine-linux-musl/bin/

ENTRYPOINT ["crystal"]

CMD ["-v"]
