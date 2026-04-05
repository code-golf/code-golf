FROM alpine:3.23 AS builder

RUN apk add --no-cache g++ linux-headers make

ENV VER=0.715

WORKDIR /luau

ADD https://github.com/luau-lang/luau.git#$VER .

RUN make config=release luau \
 && strip build/release/luau

COPY luau.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/luau -s luau.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1      /lib/
COPY --from=0 /usr/bin/luau            /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/
COPY --from=0 /luau/build/release/luau /usr/local/bin/

ENTRYPOINT ["luau"]
