FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc m4 make musl-dev

ENV VER=0.7.3

RUN curl -#L https://ftp.gnu.org/gnu/gforth/gforth-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure     \
    --prefix=/usr   \
    --without-check \
 && make            \
 && make install

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/gforth          /usr/bin/forth
COPY --from=0 /usr/include/gforth      /usr/include/gforth
COPY --from=0 /usr/lib/gforth          /usr/lib/gforth
COPY --from=0 /usr/share/gforth        /usr/share/gforth

ENTRYPOINT ["forth"]

CMD ["--version"]
