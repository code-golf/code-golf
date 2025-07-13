ARG VER=3.2.2

FROM alpine:3.22 AS builder

ARG VER

RUN apk add --no-cache curl gcc musl-dev

RUN curl -# https://downloads.freepascal.org/fpc/dist/$VER/x86_64-linux/fpc-$VER.x86_64-linux.tar \
  | tar x --strip-components 1 && sh install.sh

COPY pascal.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/pascal -s pascal.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/ld-musl-x86_64.so.1                        /lib/
COPY --from=0 /usr/bin/fpc                                    \
              /usr/bin/ld                                     \
              /usr/bin/pascal                                 \
              /usr/bin/ppcx64                                 /usr/bin/
COPY --from=0 /usr/lib/libbfd-2.44.so                         \
              /usr/lib/libctf.so.0                            \
              /usr/lib/libjansson.so.4                        \
              /usr/lib/libsframe.so.1                         \
              /usr/lib/libz.so.1                              \
              /usr/lib/libzstd.so.1                           /usr/lib/
COPY --from=0 /usr/lib/fpc/$VER/units/x86_64-linux/fcl-base   \
              /usr/lib/fpc/$VER/units/x86_64-linux/rtl        \
              /usr/lib/fpc/$VER/units/x86_64-linux/rtl-objpas /usr/lib/fpc/$VER/units/x86_64-linux/rtl/

ENTRYPOINT ["pascal"]

CMD ["--version"]
