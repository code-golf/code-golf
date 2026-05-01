FROM alpine:3.22 AS builder

ENV VER=2.49.92

RUN apk add --no-cache curl gcc make musl-dev

RUN curl -#L https://haible.de/bruno/gnu/clisp-$VER.tar.bz2 \
  | tar xj --strip-components 1

WORKDIR /src

RUN FORCE_UNSAFE_CONFIGURE=1       \
    ../configure                   \
    --ignore-absence-of-libsigsegv \
    --prefix=/usr                  \
 && make -j`nproc` install

FROM codegolf/lang-base

COPY --from=0 --parents /lib/ld-musl-*.so.1 \
                        /usr/bin/clisp      \
                        /usr/lib/clisp-*    /

ENTRYPOINT ["clisp"]

CMD ["--version"]
