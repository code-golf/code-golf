FROM alpine:3.20 AS builder

RUN apk add --no-cache npm

RUN npm --global install jsfuck@0.4.0

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh                             /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                     \
              /lib/libcrypto.so.3                          \
              /lib/libssl.so.3                             \
              /lib/libz.so.1                               /lib/
COPY --from=0 /usr/bin/env /usr/bin/node                   /usr/bin/
COPY --from=0 /usr/lib/libada.so.2                         \
              /usr/lib/libbrotlicommon.so.1                \
              /usr/lib/libbrotlidec.so.1                   \
              /usr/lib/libbrotlienc.so.1                   \
              /usr/lib/libcares.so.2                       \
              /usr/lib/libgcc_s.so.1                       \
              /usr/lib/libicudata.so.74                    \
              /usr/lib/libicui18n.so.74                    \
              /usr/lib/libicuuc.so.74                      \
              /usr/lib/libnghttp2.so.14                    \
              /usr/lib/libstdc++.so.6                      \
              /usr/lib/libuv.so.1                          /usr/lib/
COPY --from=0 /usr/local/bin                               /usr/local/bin
COPY --from=0 /usr/local/lib/node_modules/jsfuck/fuck.js   \
              /usr/local/lib/node_modules/jsfuck/jsfuck.js /usr/local/lib/node_modules/jsfuck/

COPY /jsfuckwrapper /usr/bin/

ENTRYPOINT ["jsfuckwrapper"]
