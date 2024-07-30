FROM alpine:3.20 AS builder

RUN apk add --no-cache bash curl npm

WORKDIR /usr/src

RUN curl -L https://github.com/DanielXMoore/Civet/archive/2677d5d.tar.gz | tar xz --strip=1

RUN npm i && npm run build

FROM codegolf/lang-base

# Between this comment and the next are all recorded open traces from `strace`, except the very first which I got from `ldd` wtf
COPY --from=0 /lib/ld-musl-x86_64.so.1                      \
              /lib/libcrypto.so.3                           \
              /lib/libssl.so.3                              \
              /lib/libz.so.1                                /lib/
COPY --from=0 /usr/bin/node                                 /usr/bin/
COPY --from=0 /usr/lib/libada.so.2                          \
              /usr/lib/libbrotlicommon.so.1                 \
              /usr/lib/libbrotlidec.so.1                    \
              /usr/lib/libbrotlienc.so.1                    \
              /usr/lib/libcares.so.2                        \
              /usr/lib/libgcc_s.so.1                        \
              /usr/lib/libicudata.so.74                     \
              /usr/lib/libicui18n.so.74                     \
              /usr/lib/libicuuc.so.74                       \
              /usr/lib/libnghttp2.so.14                     \
              /usr/lib/libstdc++.so.6                       \
              /usr/lib/libuv.so.1                           /usr/lib/
COPY --from=0 /usr/src/package.json                         \
              /usr/src/register.js                          /usr/local/
COPY --from=0 /usr/src/dist/civet                           \
              /usr/src/dist/config.js                       \
              /usr/src/dist/main.js                         \
              /usr/src/dist/ts-diagnostic.js                /usr/local/bin/
COPY --from=0 /usr/src/dist/unplugin/unplugin.js            /usr/local/bin/unplugin/
COPY --from=0 /usr/src/node_modules/@cspotcode              /usr/local/node_modules/@cspotcode
COPY --from=0 /usr/src/node_modules/@jridgewell             /usr/local/node_modules/@jridgewell
COPY --from=0 /usr/src/node_modules/@typescript             /usr/local/node_modules/@typescript
COPY --from=0 /usr/src/node_modules/acorn                   /usr/local/node_modules/acorn
COPY --from=0 /usr/src/node_modules/unplugin                /usr/local/node_modules/unplugin
COPY --from=0 /usr/src/node_modules/webpack-sources         /usr/local/node_modules/webpack-sources
COPY --from=0 /usr/src/node_modules/webpack-virtual-modules /usr/local/node_modules/webpack-virtual-modules
COPY --from=0 /usr/share/icu/74.2/icudt74l.dat              /usr/share/icu/74.2/
# The whole fix is probably an env thing again, but I have no F clue for now and I really want to find it solo ^^

#COPY --from=0 / /
#COPY --from=0 /usr/src /usr/local

ENTRYPOINT ["civet"]

CMD ["-v"]
