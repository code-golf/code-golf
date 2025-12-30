ARG VER=0.8.0

FROM codegolf/lang-scheme

FROM alpine:3.23 AS builder

ARG VER

RUN apk add --no-cache curl dash gcc gmp-dev make musl-dev patch

ENV SCHEME='scheme'

COPY --from=0 /boot /boot
COPY --from=0 /usr  /usr

RUN curl -#L https://www.idris-lang.org/releases/idris2-$VER.tgz \
  | tar xz --strip-components 1

COPY idris.patch /

RUN patch -p0 < idris.patch

RUN make -j`nproc` bootstrap && make install

COPY idris.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/idris -s -static idris.c

FROM codegolf/lang-base

ARG VER

COPY --from=1 /bin/uname                                      /bin/
COPY --from=1 /usr/bin/dash                                   /bin/sh
COPY --from=1 /boot/ta6le/petite.boot                         \
              /boot/ta6le/scheme.boot                         /boot/ta6le/
COPY --from=1 /lib/ld-musl-*.so.1                             /lib/
COPY --from=1 /usr/bin/dirname                                \
              /usr/bin/idris                                  \
              /usr/bin/readlink                               \
              /usr/bin/scheme                                 /usr/bin/
COPY --from=1 /usr/local/bin/idris2                           /usr/local/bin/idris
COPY --from=1 /usr/local/bin/idris2_app/idris2.so             \
              /usr/local/bin/idris2_app/libidris2_support.so  /usr/local/bin/idris2_app/
COPY --from=1 /usr/local/idris2-$VER/lib/libidris2_support.so /usr/local/idris2-$VER/lib/
COPY --from=1 /usr/local/idris2-$VER/base-$VER                /usr/local/idris2-$VER/base-$VER
COPY --from=1 /usr/local/idris2-$VER/prelude-$VER             /usr/local/idris2-$VER/prelude-$VER
COPY --from=1 /usr/local/idris2-$VER/support/chez/support.ss  /usr/local/idris2-$VER/support/chez/

ENTRYPOINT ["idris"]

CMD ["--version"]
