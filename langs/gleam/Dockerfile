FROM rust:1.83.0-alpine3.21 AS builder

RUN apk add --no-cache build-base curl icu-data npm

ENV VER=1.6.3

RUN curl -#L https://github.com/gleam-lang/gleam/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` install

FROM codegolf/lang-base

COPY --from=0 /bin                             /bin
COPY --from=0 /lib/ld-musl-x86_64.so.1         /lib/
COPY --from=0 /usr/bin/node                    \
              /usr/local/cargo/bin/gleam       /usr/bin/
COPY --from=0 /usr/lib/libada.so.2             \
              /usr/lib/libbrotlicommon.so.1    \
              /usr/lib/libbrotlidec.so.1       \
              /usr/lib/libbrotlienc.so.1       \
              /usr/lib/libcares.so.2           \
              /usr/lib/libcrypto.so.3          \
              /usr/lib/libgcc_s.so.1           \
              /usr/lib/libicudata.so.74        \
              /usr/lib/libicui18n.so.74        \
              /usr/lib/libicuuc.so.74          \
              /usr/lib/libnghttp2.so.14        \
              /usr/lib/libsimdjson.so.23       \
              /usr/lib/libsimdutf.so.11        \
              /usr/lib/libsqlite3.so.0         \
              /usr/lib/libssl.so.3             \
              /usr/lib/libstdc++.so.6          \
              /usr/lib/libz.so.1               /usr/lib/
COPY --from=0 /usr/share/icu/74.2/icudt74l.dat /usr/share/icu/

RUN chmod 777 /tmp
USER nobody
COPY --chown=nobody /gleam.toml /gleam/

RUN cd gleam \
 && gleam deps download

COPY /gleamwrapper /usr/bin/

ENTRYPOINT ["gleamwrapper"]

CMD ["--version"]
