FROM alpine:3.20 AS builder

RUN apk add --no-cache curl npm

RUN curl -L https://github.com/microsoft/TypeScript/archive/refs/tags/v5.6.2.tar.gz \
  | tar xz --strip-components 1

RUN npm install @swc/core ts-node \
 && npm run build

FROM codegolf/lang-base

COPY --from=0 /bin/busybox /bin/tsc                                                 /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                                              \
              /lib/libcrypto.so.3                                                   \
              /lib/libssl.so.3                                                      \
              /lib/libz.so.1                                                        \
              /lib/tsc.js                                                           /lib/
COPY --from=0 /node_modules/.bin                                                    /node_modules/.bin
COPY --from=0 /node_modules/@cspotcode/source-map-support/package.json              \
              /node_modules/@cspotcode/source-map-support/source-map-support.js     /node_modules/@cspotcode/source-map-support/
COPY --from=0 /node_modules/@jridgewell/resolve-uri/package.json                    /node_modules/@jridgewell/resolve-uri/
COPY --from=0 /node_modules/@jridgewell/resolve-uri/dist/resolve-uri.umd.js         /node_modules/@jridgewell/resolve-uri/dist/
COPY --from=0 /node_modules/@jridgewell/sourcemap-codec/package.json                /node_modules/@jridgewell/sourcemap-codec/
COPY --from=0 /node_modules/@jridgewell/sourcemap-codec/dist/sourcemap-codec.umd.js /node_modules/@jridgewell/sourcemap-codec/dist/
COPY --from=0 /node_modules/@jridgewell/trace-mapping/package.json                  /node_modules/@jridgewell/trace-mapping/
COPY --from=0 /node_modules/@jridgewell/trace-mapping/dist/trace-mapping.umd.js     /node_modules/@jridgewell/trace-mapping/dist/
COPY --from=0 /node_modules/@swc/core/binding.js                                    \
              /node_modules/@swc/core/index.js                                      \
              /node_modules/@swc/core/package.json                                  \
              /node_modules/@swc/core/spack.js                                      /node_modules/@swc/core/
COPY --from=0 /node_modules/@swc/core-linux-x64-musl/package.json                   \
              /node_modules/@swc/core-linux-x64-musl/swc.linux-x64-musl.node        /node_modules/@swc/core-linux-x64-musl/
COPY --from=0 /node_modules/@tsconfig/node16/*.json                                 /node_modules/@tsconfig/node16/
COPY --from=0 /node_modules/arg/index.js /node_modules/arg/package.json             /node_modules/arg/
COPY --from=0 /node_modules/make-error/index.js                                     \
              /node_modules/make-error/package.json                                 /node_modules/make-error/
COPY --from=0 /node_modules/ts-node                                                 /node_modules/ts-node
COPY --from=0 /node_modules/typescript/package.json                                 /node_modules/typescript/
COPY --from=0 /node_modules/typescript/lib/typescript.js                            /node_modules/typescript/lib/
COPY --from=0 /node_modules/yn/index.js                                             \
              /node_modules/yn/lenient.js                                           \
              /node_modules/yn/package.json                                         /node_modules/yn/
COPY --from=0 /usr/bin/env /usr/bin/ldd /usr/bin/node                               /usr/bin/
COPY --from=0 /usr/lib/libada.so.2                                                  \
              /usr/lib/libbrotlicommon.so.1                                         \
              /usr/lib/libbrotlidec.so.1                                            \
              /usr/lib/libbrotlienc.so.1                                            \
              /usr/lib/libcares.so.2                                                \
              /usr/lib/libgcc_s.so.1                                                \
              /usr/lib/libicudata.so.74                                             \
              /usr/lib/libicui18n.so.74                                             \
              /usr/lib/libicuuc.so.74                                               \
              /usr/lib/libnghttp2.so.14                                             \
              /usr/lib/libstdc++.so.6                                               \
              /usr/lib/libuv.so.1                                                   /usr/lib/

ENTRYPOINT ["tsc"]

CMD ["--version"]
