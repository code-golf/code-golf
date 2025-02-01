FROM alpine:3.21 AS builder

RUN apk add --no-cache bison cmake flex g++ gc-dev make

ENV VER=4faf069

RUN wget -O- https://github.com/Wodan58/Moy/tarball/$VER \
  | tar xz --strip-components 1

WORKDIR /build

RUN cmake ..        \
 && cmake --build . \
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
