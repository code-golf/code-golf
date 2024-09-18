FROM alpine:3.20 AS builder

RUN apk add --no-cache bash curl icu-data npm

WORKDIR /usr/src

RUN curl -L https://github.com/DanielXMoore/Civet/archive/refs/tags/v0.7.30.tar.gz \
  | tar xz --strip-components 1

RUN npm install && npm run build

FROM codegolf/lang-base

COPY --from=0 /bin/busybox                                                       /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                                           \
              /lib/libcrypto.so.3                                                \
              /lib/libssl.so.3                                                   \
              /lib/libz.so.1                                                     /lib/
COPY --from=0 /usr/bin/env /usr/bin/node                                         /usr/bin/
COPY --from=0 /usr/lib/libada.so.2                                               \
              /usr/lib/libbrotlicommon.so.1                                      \
              /usr/lib/libbrotlidec.so.1                                         \
              /usr/lib/libbrotlienc.so.1                                         \
              /usr/lib/libcares.so.2                                             \
              /usr/lib/libgcc_s.so.1                                             \
              /usr/lib/libicudata.so.74                                          \
              /usr/lib/libicui18n.so.74                                          \
              /usr/lib/libicuuc.so.74                                            \
              /usr/lib/libnghttp2.so.14                                          \
              /usr/lib/libstdc++.so.6                                            \
              /usr/lib/libuv.so.1                                                /usr/lib/
COPY --from=0 /usr/src/package.json /usr/src/register.js                         /usr/local/
COPY --from=0 /usr/src/dist/civet                                                \
              /usr/src/dist/config.js                                            \
              /usr/src/dist/main.js                                              \
              /usr/src/dist/ts-diagnostic.js                                     /usr/local/dist/
COPY --from=0 /usr/src/dist/unplugin/unplugin.js                                 /usr/local/dist/unplugin/
COPY --from=0 /usr/src/node_modules/@typescript/vfs/package.json                 /usr/local/node_modules/@typescript/vfs/
COPY --from=0 /usr/src/node_modules/@typescript/vfs/dist/index.js                \
              /usr/src/node_modules/@typescript/vfs/dist/vfs.cjs.development.js  /usr/local/node_modules/@typescript/vfs/dist/
COPY --from=0 /usr/src/node_modules/acorn/package.json                           /usr/local/node_modules/acorn/
COPY --from=0 /usr/src/node_modules/acorn/dist/acorn.js                          /usr/local/node_modules/acorn/dist/
COPY --from=0 /usr/src/node_modules/unplugin/package.json                        /usr/local/node_modules/unplugin/
COPY --from=0 /usr/src/node_modules/unplugin/dist/index.js                       /usr/local/node_modules/unplugin/dist/
COPY --from=0 /usr/src/node_modules/webpack-sources/package.json                 /usr/local/node_modules/webpack-sources/
COPY --from=0 /usr/src/node_modules/webpack-sources/lib/index.js                 /usr/local/node_modules/webpack-sources/lib/
COPY --from=0 /usr/src/node_modules/webpack-virtual-modules/package.json         /usr/local/node_modules/webpack-virtual-modules/
COPY --from=0 /usr/src/node_modules/webpack-virtual-modules/lib/index.js         \
              /usr/src/node_modules/webpack-virtual-modules/lib/virtual-stats.js /usr/local/node_modules/webpack-virtual-modules/lib/
COPY --from=0 /usr/share/icu/74.2/icudt74l.dat                                   /usr/share/icu/

ENTRYPOINT ["/usr/local/dist/civet"]

CMD ["--version"]
