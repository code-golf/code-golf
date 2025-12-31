FROM alpine:3.23 AS builder

RUN apk add --no-cache bison cmake curl flex fts-dev g++ libpcap-dev linux-headers make openssl-dev python3 zlib-dev

ENV VER=8.0.4

RUN curl -#L https://download.zeek.org/zeek-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure                      \
    --build-type='MinSizeRel'        \
    --disable-af-packet              \
    --disable-auxtools               \
    --disable-broker-tests           \
    --disable-btest                  \
    --disable-btest-pcaps            \
    --disable-cluster-backend-zeromq \
    --disable-cpp-tests              \
    --disable-javascript             \
    --disable-port-prealloc          \
    --disable-python                 \
    --disable-spicy                  \
    --disable-zeek-client            \
    --disable-zeekctl                \
    --disable-zkg                    \
    --prefix=/usr/local              \
 && make -j`nproc` install           \
 && strip /usr/local/bin/zeek        \
          /usr/local/lib/*.so.*

COPY zeek.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/zeek -s -static zeek.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1                   /lib/
COPY --from=0 /usr/bin/zeek                         /usr/bin/
COPY --from=0 /usr/lib/libcrypto.so.3               \
              /usr/lib/libfts.so.0                  \
              /usr/lib/libgcc_s.so.1                \
              /usr/lib/libpcap.so.1                 \
              /usr/lib/libssl.so.3                  \
              /usr/lib/libstdc++.so.6               \
              /usr/lib/libz.so.1                    /usr/lib/
COPY --from=0 /usr/local/bin/zeek                   /usr/local/bin/
COPY --from=0 /usr/local/lib/libbinpac.so.0         \
              /usr/local/lib/libbroker.so.4         /usr/local/lib/
COPY --from=0 /usr/local/share/zeek/base            /usr/local/share/zeek/base
COPY --from=0 /usr/local/share/zeek/builtin-plugins /usr/local/share/zeek/builtin-plugins

ENTRYPOINT ["zeek"]

CMD ["--version"]
