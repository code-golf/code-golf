FROM alpine:3.21 AS builder

RUN apk add --no-cache bison bzip2-dev cmake curl flex g++ git gmp-dev linux-headers mpfr-dev ninja openssl-dev pcre-dev zlib-dev

ENV VER=1.19.2

RUN curl -#L https://github.com/qorelanguage/qore/archive/refs/tags/release-$VER.tar.gz | tar xz

WORKDIR /qore-release-$VER/build

RUN cmake -G Ninja                \
    -DCMAKE_BUILD_TYPE='Release'  \
    -DSINGLE_COMPILATION_UNIT='1' ..

RUN ninja install

COPY qore.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/qore -s -static /qore.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1     /lib/
COPY --from=0 /usr/bin/qore                /usr/bin/
COPY --from=0 /usr/lib/libbz2.so.1         \
              /usr/lib/libcrypto.so.3      \
              /usr/lib/libgcc_s.so.1       \
              /usr/lib/libgmp.so.10        \
              /usr/lib/libmpfr.so.6        \
              /usr/lib/libpcre.so.1        \
              /usr/lib/libssl.so.3         \
              /usr/lib/libstdc++.so.6      \
              /usr/lib/libz.so.1           /usr/lib/
COPY --from=0 /usr/local/bin/qore          /usr/local/bin/
COPY --from=0 /usr/local/lib/libqore.so.12 /usr/local/lib/

ENTRYPOINT ["qore"]

CMD ["--version"]
