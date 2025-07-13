FROM alpine:3.22 AS builder

RUN apk add --no-cache clang cmake curl make patch

ENV VER=3.2

RUN curl -#L https://github.com/albertodemichelis/squirrel/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

COPY argv.patch squirrel.c /

RUN patch -p0 < argv.patch

WORKDIR /build

RUN cmake .. \
 && make install

RUN strip /usr/local/lib/libsqstdlib.so.0 \
          /usr/local/lib/libsquirrel.so.0

RUN gcc -Wall -Werror -Wextra -o /usr/bin/squirrel -s /squirrel.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1        /lib/
COPY --from=0 /usr/bin/squirrel               \
              /usr/local/bin/sq               /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1          \
              /usr/lib/libstdc++.so.6         \
              /usr/local/lib/libsqstdlib.so.0 \
              /usr/local/lib/libsquirrel.so.0 /usr/lib/

ENTRYPOINT ["squirrel"]

CMD ["--version"]
