FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc linux-headers make musl-dev

ENV LDFLAGS='-static' VER=4.10

ADD --unpack https://ftpmirror.gnu.org/gnu/sed/sed-$VER.tar.xz /

WORKDIR /sed-$VER

RUN ./configure   \
    --enable-lto  \
    --prefix=/usr \
 && make install  \
 && strip /usr/bin/sed

FROM codegolf/lang-base

COPY --from=0 /usr/bin/sed /usr/bin/

ENTRYPOINT ["sed"]

CMD ["--version"]
