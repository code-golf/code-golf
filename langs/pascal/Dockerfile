FROM alpine:3.23 AS builder

RUN apk add --no-cache curl gcc musl-dev

ENV VER=3.2.2

RUN curl -# https://downloads.freepascal.org/fpc/dist/$VER/x86_64-linux/fpc-$VER.x86_64-linux.tar \
  | tar x --strip-components 1 && sh install.sh

COPY pascal.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/pascal -s pascal.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1                            /lib/
COPY --from=0 /usr/bin/fpc                                   \
              /usr/bin/ld                                    \
              /usr/bin/pascal                                \
              /usr/bin/ppcx64                                /usr/bin/
COPY --from=0 /usr/lib/libbfd-2.45.1.so                      \
              /usr/lib/libctf.so.0                           \
              /usr/lib/libjansson.so.4                       \
              /usr/lib/libsframe.so.2                        \
              /usr/lib/libz.so.1                             \
              /usr/lib/libzstd.so.1                          /usr/lib/
COPY --parents                                               \
     --from=0 /usr/lib/./fpc/*/units/x86_64-linux/fcl-base   \
              /usr/lib/./fpc/*/units/x86_64-linux/rtl        \
              /usr/lib/./fpc/*/units/x86_64-linux/rtl-objpas /usr/lib/

ENTRYPOINT ["pascal"]

CMD ["--version"]
