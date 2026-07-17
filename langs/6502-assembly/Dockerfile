FROM alpine:3.23 AS builder

RUN apk add --no-cache npm

ENV VER=0.2.0

RUN npm install --prefix=/usr fluffy-6502@$VER

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1           /lib/
COPY --from=0 /usr/bin/node                 /usr/bin/
COPY --from=0 /usr/lib/libada.so.3          \
              /usr/lib/libbrotlicommon.so.1 \
              /usr/lib/libbrotlidec.so.1    \
              /usr/lib/libbrotlienc.so.1    \
              /usr/lib/libcares.so.2        \
              /usr/lib/libcrypto.so.3       \
              /usr/lib/libgcc_s.so.1        \
              /usr/lib/libicudata.so.76     \
              /usr/lib/libicui18n.so.76     \
              /usr/lib/libicuuc.so.76       \
              /usr/lib/libnghttp2.so.14     \
              /usr/lib/libsimdjson.so.25    \
              /usr/lib/libsimdutf.so.27     \
              /usr/lib/libsqlite3.so.0      \
              /usr/lib/libssl.so.3          \
              /usr/lib/libstdc++.so.6       \
              /usr/lib/libz.so.1            \
              /usr/lib/libzstd.so.1         /usr/lib/
COPY --from=0 /usr/node_modules/fluffy-6502 /usr/node_modules/fluffy-6502

COPY --chmod=755 runner.js /usr/node_modules/fluffy-6502
