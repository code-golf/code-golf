ARG VER=2.49.92

FROM alpine:3.22 AS builder

ARG VER

RUN apk add --no-cache build-base curl

RUN curl -#L https://haible.de/bruno/gnu/clisp-$VER.tar.bz2 \
  | tar xj --strip-components 1

WORKDIR /src

RUN FORCE_UNSAFE_CONFIGURE=1       \
    ../configure                   \
    --ignore-absence-of-libsigsegv \
    --prefix=/usr                  \
 && make install

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/clisp           /usr/bin/common-lisp
COPY --from=0 /usr/lib/clisp-$VER      /usr/lib/clisp-$VER

ENTRYPOINT ["common-lisp"]

CMD ["--version"]
