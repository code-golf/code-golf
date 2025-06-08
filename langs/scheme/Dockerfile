FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl

ENV VER=10.2.0

RUN curl -#L https://github.com/cisco/ChezScheme/releases/download/v$VER/csv$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure      \
    --disable-curses \
    --disable-x11    \
 && make -j`nproc`   \
 && make install

RUN strip /usr/bin/scheme

COPY scheme.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/schemewrapper -s scheme.c

FROM codegolf/lang-base

COPY --from=0 /boot/pb/petite.boot     \
              /boot/pb/scheme.boot     /boot/pb/
COPY --from=0 /boot/ta6le/petite.boot  \
              /boot/ta6le/scheme.boot  /boot/ta6le/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/scheme          \
              /usr/bin/schemewrapper   /usr/bin/

ENTRYPOINT ["schemewrapper"]

CMD ["--version"]
