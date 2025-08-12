FROM alpine:3.22 AS builder

RUN apk add --no-cache curl clang musl-dev make git

ENV VER=73e9219 COMPILER=clang STANDARD=EXTENSION

WORKDIR /usr/bin

RUN curl -#L https://github.com/synt7x/knightjit/tarball/$VER \
  | tar xz --strip-components 1

RUN make CFLAGS='-O3 -march=x86-64-v3 -Wall -Wextra -std=c99' \
 && strip -s /usr/bin/target/knight

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/target/knight /usr/bin/

ENTRYPOINT ["knight"]
