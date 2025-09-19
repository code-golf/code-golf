FROM alpine:3.22 AS builder

RUN apk add --no-cache curl g++ linux-headers make

ENV VER=0.691

RUN curl -#L https://github.com/luau-lang/luau/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

RUN make config=release luau \
 && strip build/release/luau

COPY luau.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/luau -s luau.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/luau            /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/
COPY --from=0 /build/release/luau      /usr/local/bin/

ENTRYPOINT ["luau"]
