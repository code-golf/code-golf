FROM alpine:3.23 AS builder

RUN apk add --no-cache curl gcc linux-headers make musl-dev

ENV LDFLAGS='-static' VER=4.10

RUN curl -#L https://ftpmirror.gnu.org/gnu/sed/sed-$VER.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure   \
    --enable-lto  \
    --prefix=/usr \
 && make install  \
 && strip /usr/bin/sed

FROM codegolf/lang-base

COPY --from=0 /usr/bin/sed /usr/bin/

ENTRYPOINT ["sed"]

CMD ["--version"]
