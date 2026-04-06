FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV VER=2.49.92

RUN curl -#L https://haible.de/bruno/gnu/clisp-$VER.tar.bz2 \
  | tar xj --strip-components 1

WORKDIR /src

RUN FORCE_UNSAFE_CONFIGURE=1       \
    ../configure                   \
    --ignore-absence-of-libsigsegv \
    --prefix=/usr                  \
 && make -j`nproc` install

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/clisp      /usr/bin/common-lisp
COPY --parents                    \
     --from=0 /usr/lib/./clisp-*  /usr/lib/

ENTRYPOINT ["common-lisp"]

CMD ["--version"]
