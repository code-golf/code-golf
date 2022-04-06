FROM crystallang/crystal:1.3.2-alpine as builder

RUN strip /usr/bin/crystal

FROM codegolf/lang-base

ENV CRYSTAL_CACHE_DIR=/tmp PATH=/usr/bin:/bin

COPY --from=0 /bin/sh                       /bin/
COPY --from=0 /lib                          /lib
COPY --from=0 /usr/bin/cc /usr/bin/crystal  /usr/bin/
COPY --from=0 /usr/include                  /usr/include
COPY --from=0 /usr/lib                      /usr/lib
COPY --from=0 /usr/libexec                  /usr/libexec
COPY --from=0 /usr/local                    /usr/local
COPY --from=0 /usr/share                    /usr/share
COPY --from=0 /usr/x86_64-alpine-linux-musl /usr/x86_64-alpine-linux-musl

ENTRYPOINT ["crystal"]

CMD ["-v"]
