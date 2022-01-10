FROM crystallang/crystal:1.3.0-alpine as builder

# Crystal expects /dev/null to exist.
RUN mkdir /empty /mydev && mknod -m 0666 /mydev/null c 1 3

RUN strip /usr/bin/crystal

FROM scratch

ENV CRYSTAL_CACHE_DIR=/tmp PATH=/usr/bin:/bin

COPY --from=0 /bin/sh                       /bin/
COPY --from=0 /mydev                        /dev
COPY --from=0 /lib                          /lib
COPY --from=0 /empty                        /proc
COPY --from=0 /empty                        /tmp
COPY --from=0 /usr/bin/cc /usr/bin/crystal  /usr/bin/
COPY --from=0 /usr/include                  /usr/include
COPY --from=0 /usr/lib                      /usr/lib
COPY --from=0 /usr/libexec                  /usr/libexec
COPY --from=0 /usr/local                    /usr/local
COPY --from=0 /usr/share                    /usr/share
COPY --from=0 /usr/x86_64-alpine-linux-musl /usr/x86_64-alpine-linux-musl

ENTRYPOINT ["crystal"]

CMD ["-v"]
