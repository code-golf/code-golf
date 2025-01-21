FROM alpine:3.21 AS builder

RUN apk add --no-cache cmake g++ gc-dev make

ENV VER=d7130c7

RUN wget -O- https://github.com/Wodan58/joy1/tarball/$VER \
  | tar xz --strip-components 1

WORKDIR /build

RUN cmake -G 'Unix Makefiles' .. \
 && cmake --build .              \
 && strip joy

COPY joywrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/joywrapper -s /joywrapper.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /build/joy               \
              /usr/bin/joywrapper      /usr/bin/
COPY --from=0 /usr/lib/libgc.so.1      /usr/lib/

ENTRYPOINT ["joywrapper"]

CMD ["--version"]
