FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl

ENV ARCH=x86_64-linux VER=3.2.2

RUN curl -# https://downloads.freepascal.org/fpc/dist/$VER/$ARCH/fpc-$VER.$ARCH.tar \
  | tar x --strip-components 1 && sh install.sh

# Remove some unneeded units.
RUN find /usr/lib/fpc/$VER/units/$ARCH \
    -type d -mindepth 1 -maxdepth 1    \
    -not -name fcl-base                \
    -not -name rtl                     \
    -not -name rtl-console             \
    -not -name rtl-objpas              \
    -exec rm -r {} +

COPY pascal.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/pascal -s pascal.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /usr/bin/fpc              \
              /usr/bin/ld               \
              /usr/bin/pascal           \
              /usr/bin/ppcx64           /usr/bin/
COPY --from=0 /usr/lib/libbfd-2.43.1.so \
              /usr/lib/libctf.so.0      \
              /usr/lib/libjansson.so.4  \
              /usr/lib/libsframe.so.1   \
              /usr/lib/libz.so.1        \
              /usr/lib/libzstd.so.1     /usr/lib/
COPY --from=0 /usr/lib/fpc              /usr/lib/fpc

ENTRYPOINT ["pascal"]

CMD ["--version"]
