FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl

ENV LDFLAGS='-static' VER=4.9

RUN curl -#L https://ftp.gnu.org/gnu/sed/sed-$VER.tar.xz \
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
