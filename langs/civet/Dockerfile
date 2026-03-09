FROM alpine:3.23 AS builder

RUN apk add --no-cache bash curl icu-data npm

ENV VER=0.11.4

RUN curl -#L https://github.com/DanielXMoore/Civet/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN npm install \
 && npm run build

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1                                        /lib/
COPY --from=0 /usr/bin/env /usr/bin/node                                 /usr/bin/
COPY --from=0 /usr/lib/libada.so.3                                       \
              /usr/lib/libbrotlicommon.so.1                              \
              /usr/lib/libbrotlidec.so.1                                 \
              /usr/lib/libbrotlienc.so.1                                 \
              /usr/lib/libcares.so.2                                     \
              /usr/lib/libcrypto.so.3                                    \
              /usr/lib/libgcc_s.so.1                                     \
              /usr/lib/libicudata.so.76                                  \
              /usr/lib/libicui18n.so.76                                  \
              /usr/lib/libicuuc.so.76                                    \
              /usr/lib/libnghttp2.so.14                                  \
              /usr/lib/libsimdjson.so.25                                 \
              /usr/lib/libsimdutf.so.27                                  \
              /usr/lib/libsqlite3.so.0                                   \
              /usr/lib/libssl.so.3                                       \
              /usr/lib/libstdc++.so.6                                    \
              /usr/lib/libz.so.1                                         \
              /usr/lib/libzstd.so.1                                      /usr/lib/
COPY --parents                                                           \
     --from=0 /package.json                                              \
              /register.js                                               \
              /dist/civet                                                \
              /dist/config.js                                            \
              /dist/main.js                                              \
              /dist/ts-diagnostic.js                                     \
              /dist/unplugin/unplugin.js                                 \
              /node_modules/@typescript/vfs/package.json                 \
              /node_modules/@typescript/vfs/dist/index.js                \
              /node_modules/@typescript/vfs/dist/vfs.cjs.development.js  \
              /node_modules/acorn/package.json                           \
              /node_modules/acorn/dist/acorn.js                          \
              /node_modules/picomatch/index.js                           \
              /node_modules/picomatch/lib/constants.js                   \
              /node_modules/picomatch/lib/parse.js                       \
              /node_modules/picomatch/lib/picomatch.js                   \
              /node_modules/picomatch/lib/scan.js                        \
              /node_modules/picomatch/lib/utils.js                       \
              /node_modules/unplugin/package.json                        \
              /node_modules/unplugin/dist/context-*.cjs                  \
              /node_modules/unplugin/dist/index.js                       \
              /node_modules/unplugin/dist/index.cjs                      \
              /node_modules/unplugin/dist/utils-*.cjs                    \
              /node_modules/unplugin/dist/webpack-like-*.cjs             \
              /node_modules/webpack-virtual-modules/package.json         \
              /node_modules/webpack-virtual-modules/lib/index.js         \
              /node_modules/webpack-virtual-modules/lib/virtual-stats.js /usr/local/
COPY --from=0 /usr/share/icu/76.1/icudt76l.dat                           /usr/share/icu/

ENTRYPOINT ["/usr/local/dist/civet"]

CMD ["--version"]
