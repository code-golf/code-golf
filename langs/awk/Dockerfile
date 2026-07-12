FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev

ENV VER=5.4.1

ADD --unpack https://ftpmirror.gnu.org/gnu/gawk/gawk-$VER.tar.xz /

WORKDIR /gawk-$VER

RUN ./configure                       \
    --disable-lint                    \
    --disable-pma                     \
    --prefix=/usr                     \
    --program-transform-name='s/^g//' \
 && make install                      \
 && strip /usr/bin/awk

FROM codegolf/lang-base

COPY --from=0 --parents /lib/ld-musl-*.so.1 \
                        /usr/bin/awk        \
                        /usr/lib/gawk       /

ENTRYPOINT ["awk"]

CMD ["--version"]
