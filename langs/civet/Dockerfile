FROM alpine:3.21 AS builder

RUN apk add --no-cache icu-data npm

ENV VER=0.9.7

RUN npm install --global --prefix=/usr @danielx/civet@$VER

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1                                                                       /lib/
COPY --from=0 /usr/bin/env /usr/bin/node                                                                     /usr/bin/
COPY --from=0 /usr/lib/libada.so.2                                                                           \
              /usr/lib/libbrotlicommon.so.1                                                                  \
              /usr/lib/libbrotlidec.so.1                                                                     \
              /usr/lib/libbrotlienc.so.1                                                                     \
              /usr/lib/libcares.so.2                                                                         \
              /usr/lib/libcrypto.so.3                                                                        \
              /usr/lib/libgcc_s.so.1                                                                         \
              /usr/lib/libicudata.so.74                                                                      \
              /usr/lib/libicui18n.so.74                                                                      \
              /usr/lib/libicuuc.so.74                                                                        \
              /usr/lib/libnghttp2.so.14                                                                      \
              /usr/lib/libsimdjson.so.23                                                                     \
              /usr/lib/libsimdutf.so.11                                                                      \
              /usr/lib/libsqlite3.so.0                                                                       \
              /usr/lib/libssl.so.3                                                                           \
              /usr/lib/libstdc++.so.6                                                                        \
              /usr/lib/libz.so.1                                                                             /usr/lib/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/package.json                                              \
              /usr/lib/node_modules/@danielx/civet/register.js                                               /usr/lib/node_modules/@danielx/civet/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/dist/civet                                                \
              /usr/lib/node_modules/@danielx/civet/dist/config.js                                            \
              /usr/lib/node_modules/@danielx/civet/dist/main.js                                              \
              /usr/lib/node_modules/@danielx/civet/dist/ts-diagnostic.js                                     /usr/lib/node_modules/@danielx/civet/dist/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/dist/unplugin/unplugin.js                                 /usr/lib/node_modules/@danielx/civet/dist/unplugin/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/node_modules/@typescript/vfs/package.json                 /usr/lib/node_modules/@typescript/vfs/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/node_modules/@typescript/vfs/dist/index.js                \
              /usr/lib/node_modules/@danielx/civet/node_modules/@typescript/vfs/dist/vfs.cjs.development.js  /usr/lib/node_modules/@typescript/vfs/dist/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/node_modules/acorn/package.json                           /usr/lib/node_modules/acorn/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/node_modules/acorn/dist/acorn.js                          /usr/lib/node_modules/acorn/dist/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/node_modules/unplugin/package.json                        /usr/lib/node_modules/unplugin/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/node_modules/unplugin/dist/context-CYCvfJ5z.cjs           \
              /usr/lib/node_modules/@danielx/civet/node_modules/unplugin/dist/context-C3WvS-qG.cjs           \
              /usr/lib/node_modules/@danielx/civet/node_modules/unplugin/dist/context-D3o27x9-.cjs           \
              /usr/lib/node_modules/@danielx/civet/node_modules/unplugin/dist/utils-CV7Znnzx.cjs             \
              /usr/lib/node_modules/@danielx/civet/node_modules/unplugin/dist/index.js                       \
              /usr/lib/node_modules/@danielx/civet/node_modules/unplugin/dist/index.cjs                      \
              /usr/lib/node_modules/@danielx/civet/node_modules/unplugin/dist/webpack-like-Cz79ljvC.cjs      /usr/lib/node_modules/unplugin/dist/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/node_modules/webpack-virtual-modules/package.json         /usr/lib/node_modules/webpack-virtual-modules/
COPY --from=0 /usr/lib/node_modules/@danielx/civet/node_modules/webpack-virtual-modules/lib/index.js         \
              /usr/lib/node_modules/@danielx/civet/node_modules/webpack-virtual-modules/lib/virtual-stats.js /usr/lib/node_modules/webpack-virtual-modules/lib/
COPY --from=0 /usr/share/icu/74.2/icudt74l.dat                                                               /usr/share/icu/

ENTRYPOINT ["/usr/lib/node_modules/@danielx/civet/dist/civet"]

CMD ["--version"]
