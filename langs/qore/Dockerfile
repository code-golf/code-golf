FROM alpine:3.22 AS builder

RUN apk add --no-cache bison build-base bzip2-dev curl flex gmp-dev \
                       linux-headers mpfr-dev openssl-dev pcre2-dev zlib-dev

ENV VER=2.2.0

RUN curl -#L https://github.com/qorelanguage/qore/releases/download/release-$VER/qore-$VER.tar.bz2 \
  | tar xj --strip-components 1

RUN sed -i 's/exit(0);//' command-line.cpp

RUN ./configure            \
    --disable-debug        \
    --prefix=/usr          \
 && make -j`nproc` install \
 && strip /usr/bin/qore    \
          /usr/lib/libqore.so.12

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1      /lib/
COPY --from=0 /usr/bin/qore            /usr/bin/
COPY --from=0 /usr/lib/libbz2.so.1     \
              /usr/lib/libcrypto.so.3  \
              /usr/lib/libgcc_s.so.1   \
              /usr/lib/libgmp.so.10    \
              /usr/lib/libmpfr.so.6    \
              /usr/lib/libpcre2-8.so.0 \
              /usr/lib/libqore.so.12   \
              /usr/lib/libssl.so.3     \
              /usr/lib/libstdc++.so.6  \
              /usr/lib/libz.so.1       /usr/lib/

ENTRYPOINT ["qore"]

CMD ["--version"]
