ARG VER=0.7.0

FROM alpine:3.21 AS builder

ARG VER

RUN apk add --no-cache build-base curl gmp-dev

RUN curl -#L https://github.com/idris-lang/Idris2/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

COPY config.patch /

RUN patch -p0 < config.patch

COPY --from=codegolf/lang-scheme . /

# Idris 2 has two major building stages: the first stage builds the compiler,
# the second stage uses the built compiler to rebuild (parts of) the project.
RUN SCHEME='scheme' make -j`nproc` bootstrap && make install

COPY idris.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/idris -s -static idris.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /bin/sh /bin/uname                              /bin/
COPY --from=0 /boot/ta6le/petite.boot                         \
              /boot/ta6le/scheme.boot                         /boot/ta6le/
COPY --from=0 /lib/ld-musl-x86_64.so.1                        /lib/
COPY --from=0 /usr/bin/dirname                                \
              /usr/bin/idris                                  \
              /usr/bin/readlink                               \
              /usr/bin/scheme                                 /usr/bin/
COPY --from=0 /usr/local/bin/idris2                           /usr/local/bin/idris
COPY --from=0 /usr/local/bin/idris2_app/idris2.so             \
              /usr/local/bin/idris2_app/libidris2_support.so  /usr/local/bin/idris2_app/
COPY --from=0 /usr/local/idris2-$VER/lib/libidris2_support.so /usr/local/idris2-$VER/lib/
COPY --from=0 /usr/local/idris2-$VER/base-$VER                /usr/local/idris2-$VER/base-$VER
COPY --from=0 /usr/local/idris2-$VER/prelude-$VER             /usr/local/idris2-$VER/prelude-$VER
COPY --from=0 /usr/local/idris2-$VER/support/chez/support.ss  /usr/local/idris2-$VER/support/chez/

ENTRYPOINT ["idris"]

CMD ["--version"]
