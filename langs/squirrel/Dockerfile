FROM alpine:3.21 AS builder

RUN apk add --no-cache clang cmake curl make patch

ENV VER=3.2

RUN curl -#L https://github.com/albertodemichelis/squirrel/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

COPY argv.patch \
     squirrel.c /

RUN patch -p 0 < argv.patch

RUN gcc -Wall -Werror -Wextra -o /usr/bin/squirrel -s -static squirrel.c

WORKDIR /build

RUN cmake -DCMAKE_INSTALL_PREFIX=/usr .. \
 && make install

RUN strip /usr/lib/libsqstdlib.so.0 \
          /usr/lib/libsquirrel.so.0

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/bin/sq               \
              /usr/bin/squirrel         /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1    \
              /usr/lib/libsqstdlib.so.0 \
              /usr/lib/libsquirrel.so.0 \
              /usr/lib/libstdc++.so.6   /usr/lib/

ENTRYPOINT ["squirrel"]

CMD ["--version"]
