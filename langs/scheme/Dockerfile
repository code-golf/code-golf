FROM alpine:3.23 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV VER=10.4.1

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

COPY --from=0           /usr/bin/scheme        /usr/bin/
COPY --from=0 --parents /boot/*/*.boot         \
                        /lib/ld-musl-*.so.1    \
                        /usr/bin/schemewrapper /

ENTRYPOINT ["schemewrapper"]

CMD ["--version"]
