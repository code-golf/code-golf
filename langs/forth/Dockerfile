FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc m4 make musl-dev

ENV LDFLAGS='-static' VER=0.7.3

RUN curl -#L https://ftpmirror.gnu.org/gnu/gforth/gforth-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure     \
    --prefix=/usr   \
    --without-check \
 && make            \
 && make install    \
 && strip /usr/bin/gforth

FROM codegolf/lang-base

COPY --from=0 /usr/bin/gforth             /usr/bin/forth
COPY --from=0 /usr/lib/gforth/*/gforth.fi /usr/lib/gforth/site-forth/

ENTRYPOINT ["forth"]

CMD ["--version"]
