FROM alpine:3.24 AS builder

RUN apk add --no-cache curl icu-data npm

ENV VER=1.2.1

WORKDIR /stax

ADD https://github.com/tomtheisen/stax.git#v$VER .

WORKDIR /stax/StaxLang.JS

RUN npm install \
 # Remove debug console.log that gets captured by the cli wrapper
 && sed -i '/console\.log("macro", A2S(a))/d' stax.ts \
 && npx tsc

COPY stax.js /usr/bin/stax

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1              /lib/
COPY --from=0 /usr/bin/node                    \
              /usr/bin/stax                    /usr/bin/
COPY --from=0 /usr/lib/libada.so.3             \
              /usr/lib/libbrotlicommon.so.1    \
              /usr/lib/libbrotlidec.so.1       \
              /usr/lib/libbrotlienc.so.1       \
              /usr/lib/libcares.so.2           \
              /usr/lib/libcrypto.so.3          \
              /usr/lib/libgcc_s.so.1           \
              /usr/lib/libicudata.so.78        \
              /usr/lib/libicui18n.so.78        \
              /usr/lib/libicuuc.so.78          \
              /usr/lib/libnghttp2.so.14        \
              /usr/lib/libsimdjson.so.29       \
              /usr/lib/libsimdutf.so.34        \
              /usr/lib/libsqlite3.so.0         \
              /usr/lib/libssl.so.3             \
              /usr/lib/libstdc++.so.6          \
              /usr/lib/libz.so.1               \
              /usr/lib/libzstd.so.1            /usr/lib/
COPY --from=0 /stax/StaxLang.JS/js/*.js        /usr/lib/stax/
COPY --from=0 /usr/share/icu/78.1/icudt78l.dat /usr/share/icu/

ENTRYPOINT ["stax"]

CMD ["--version"]
