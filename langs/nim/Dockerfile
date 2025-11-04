FROM alpine:3.22 AS builder

RUN apk add --no-cache curl dash gcc musl-dev pcre tcc tcc-libs-static

# Also bump Arturo.
ENV VER=2.2.6

RUN curl -#L https://nim-lang.org/download/nim-$VER.tar.xz \
  | tar xJ --strip-components 1

RUN ./build.sh             \
 && ./bin/nim c koch       \
 && ./koch boot -d:release \
 && ./install.sh /         \
 && strip /nim/bin/nim

# Docker can only copy symlinks by copying the whole folder.
RUN mkdir /usr-lib-symlinks \
 && cp -P /usr/lib/libc.so /usr-lib-symlinks

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash          /bin/sh
COPY --from=0 /lib/ld-musl-*.so.1    /lib/
COPY --from=0 /nim/bin/nim           \
              /usr/bin/tcc           /usr/bin/
COPY --from=0 /usr/include           /usr/include
COPY --from=0 /usr-lib-symlinks      /usr/lib
COPY --from=0 /usr/lib/crt1.o        \
              /usr/lib/crti.o        \
              /usr/lib/crtn.o        \
              /usr/lib/libdl.a       \
              /usr/lib/libm.a        \
              /usr/lib/libpcre.so.1  \
              /usr/lib/librt.a       \
              /usr/lib/libtcc.so     /usr/lib/
COPY --from=0 /nim/lib               /usr/lib/nim/lib
COPY --from=0 /usr/lib/tcc/libtcc1.a /usr/lib/tcc/

COPY nim.cfg /etc/nim/

ENTRYPOINT ["nim"]

CMD ["--version"]
