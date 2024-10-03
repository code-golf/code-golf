FROM rust:1.81.0-alpine3.20 AS builder

RUN apk add --no-cache build-base curl icu-data npm

RUN curl -L https://github.com/gleam-lang/gleam/archive/refs/tags/v1.5.1.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` \
 && make install

FROM codegolf/lang-base

COPY --from=0 /bin                             /bin
COPY --from=0 /lib/ld-musl-x86_64.so.1         \
              /lib/libcrypto.so.3              \
              /lib/libssl.so.3                 \
              /lib/libz.so.1                   /lib/
COPY --from=0 /usr/bin/node                    \
              /usr/local/cargo/bin/gleam       /usr/bin/
COPY --from=0 /usr/lib/libada.so.2             \
              /usr/lib/libbrotlicommon.so.1    \
              /usr/lib/libbrotlidec.so.1       \
              /usr/lib/libbrotlienc.so.1       \
              /usr/lib/libcares.so.2           \
              /usr/lib/libgcc_s.so.1           \
              /usr/lib/libicudata.so.74        \
              /usr/lib/libicui18n.so.74        \
              /usr/lib/libicuuc.so.74          \
              /usr/lib/libnghttp2.so.14        \
              /usr/lib/libstdc++.so.6          \
              /usr/lib/libuv.so.1              /usr/lib/
COPY --from=0 /usr/share/icu/74.2/icudt74l.dat /usr/share/icu/

RUN chmod 777 /tmp
USER nobody
COPY --chown=nobody /gleam.toml /gleam/

RUN cd gleam \
 && gleam deps download

COPY /gleamwrapper /usr/bin/

ENTRYPOINT ["gleamwrapper"]

CMD ["--version"]
