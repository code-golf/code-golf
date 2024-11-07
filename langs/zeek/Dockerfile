FROM alpine:3.20 AS builder

RUN apk add --no-cache bison build-base cmake curl flex fts-dev libpcap-dev linux-headers openssl-dev python3 zlib-dev

RUN curl https://download.zeek.org/zeek-7.0.3.tar.gz \
  | tar xz --strip-components 1

RUN ./configure             \
    --disable-af-packet     \
    --disable-auxtools      \
    --disable-broker-tests  \
    --disable-btest         \
    --disable-btest-pcaps   \
    --disable-cpp-tests     \
    --disable-javascript    \
    --disable-port-prealloc \
    --disable-python        \
    --disable-spicy         \
    --disable-zeek-client   \
    --disable-zeekctl       \
    --disable-zkg           \
    --prefix=/usr           \
 && make -j`nproc` install

RUN strip /usr/bin/zeek

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1        \
              /lib/libz.so.1                  /lib/
COPY --from=0 /usr/bin/zeek                   /usr/bin/
COPY --from=0 /usr/lib/libbinpac.so.0         \
              /usr/lib/libbroker.so.4         \
              /usr/lib/libcrypto.so.3         \
              /usr/lib/libfts.so.0            \
              /usr/lib/libgcc_s.so.1          \
              /usr/lib/libpcap.so.1           \
              /usr/lib/libssl.so.3            \
              /usr/lib/libstdc++.so.6         /usr/lib/
COPY --from=0 /usr/share/zeek/base            /usr/share/zeek/base
COPY --from=0 /usr/share/zeek/builtin-plugins /usr/share/zeek/builtin-plugins

ENTRYPOINT ["zeek"]

CMD ["--version"]
