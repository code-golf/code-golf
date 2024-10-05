FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl gd-dev giflib-dev groff

RUN curl https://www.bertnase.de/npiet/npiet-1.3f.tar.gz | tar xz

RUN cd npiet-1.3f   \
 && ./configure     \
    --prefix="/usr" \
 && make            \
 && make install

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh              /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1      \
              /lib/libz.so.1                /lib/
COPY --from=0 /usr/bin/npiet                /usr/bin/piet
COPY --from=0 /usr/lib/libaom.so.3          \
              /usr/lib/libavif.so.16        \
              /usr/lib/libbrotlidec.so.1    \
              /usr/lib/libbrotlicommon.so.1 \
              /usr/lib/libbsd.so.0          \
              /usr/lib/libbz2.so.1          \
              /usr/lib/libdav1d.so.7        \
              /usr/lib/libexpat.so.1        \
              /usr/lib/libfontconfig.so.1   \
              /usr/lib/libfreetype.so.6     \
              /usr/lib/libgd.so.3           \
              /usr/lib/libgif.so.7          \
              /usr/lib/libjpeg.so.8         \
              /usr/lib/libmd.so.0           \
              /usr/lib/libpng16.so.16       \
              /usr/lib/libsharpyuv.so.0     \
              /usr/lib/libtiff.so.6         \
              /usr/lib/libwebp.so.7         \
              /usr/lib/libX11.so.6          \
              /usr/lib/libXau.so.6          \
              /usr/lib/libxcb.so.1          \
              /usr/lib/libXdmcp.so.6        \
              /usr/lib/libXpm.so.4          \
              /usr/lib/libzstd.so.1         /usr/lib/

COPY /pietwrapper /usr/bin/

ENTRYPOINT ["pietwrapper"]

CMD ["--version"]
