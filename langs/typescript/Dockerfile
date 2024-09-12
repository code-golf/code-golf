FROM alpine:3.20 AS builder

RUN apk add --no-cache curl npm

RUN curl -L https://github.com/microsoft/TypeScript/archive/refs/tags/v5.6.2.tar.gz \
  | tar xz --strip-components 1

RUN npm install @swc/core ts-node \
 && npm run build

FROM codegolf/lang-base

COPY --from=0 /bin/busybox /bin/tsc         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1      \
              /lib/libcrypto.so.3           \
              /lib/libssl.so.3              \
              /lib/libz.so.1                \
              /lib/tsc.js                   /lib/
COPY --from=0 /node_modules                 /node_modules
COPY --from=0 /usr/bin/env                  \
              /usr/bin/ldd                  \
              /usr/bin/node                 /usr/bin/
COPY --from=0 /usr/lib/libada.so.2          \
              /usr/lib/libbrotlicommon.so.1 \
              /usr/lib/libbrotlidec.so.1    \
              /usr/lib/libbrotlienc.so.1    \
              /usr/lib/libcares.so.2        \
              /usr/lib/libgcc_s.so.1        \
              /usr/lib/libicudata.so.74     \
              /usr/lib/libicui18n.so.74     \
              /usr/lib/libicuuc.so.74       \
              /usr/lib/libnghttp2.so.14     \
              /usr/lib/libstdc++.so.6       \
              /usr/lib/libuv.so.1           /usr/lib/

ENTRYPOINT ["tsc"]

CMD ["--version"]
