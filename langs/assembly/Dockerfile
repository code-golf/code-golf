FROM alpine:3.23 AS builder

RUN apk add --no-cache npm

ENV VER=2.4.5

RUN npm install --prefix=/usr @defasm/cli@1.2.3 @defasm/core@$VER

# Rewrite shebang to bypass /usr/bin/env and go directly to /usr/bin/node.
RUN sed -i '1c#!/usr/bin/node' /usr/node_modules/@defasm/cli/cli.js

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
COPY --from=0 /usr/node_modules/@defasm     /usr/node_modules/@defasm
