FROM alpine:3.22 AS builder

RUN apk add --no-cache clang curl git make

ENV CFLAGS='-O3 -march=x86-64-v3 -DJIT_OFF' VER=d118d76

WORKDIR /knight

RUN curl -#L https://github.com/synt7x/knightjit/tarball/$VER \
  | tar xz --strip-components 1

RUN make

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1   /lib/
COPY --from=0 /knight/target/knight /usr/bin/

ENTRYPOINT ["knight"]
