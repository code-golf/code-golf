FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl gd-dev giflib-dev git groff libffi-dev openssl-dev zlib-dev

ENV PYVER=3.13.0 VERSION=1.3f

RUN curl https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

RUN cd Python-$PYVER           \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.13  \
          /usr/lib/python3.13/lib-dynload/*.so

# Remove some unneeded libraries.
RUN cd /usr/lib/python3.13      \
 && rm -r __pycache__/doctest.* \
          __pycache__/pydoc.*   \
          __pycache__/turtle.*  \
          ensurepip             \
          idlelib               \
          pydoc_data            \
          tkinter               \
          turtle.py             \
          turtledemo

RUN pip3 install pillow

RUN curl https://www.bertnase.de/npiet/npiet-$VERSION.tar.gz \
  | tar xz --strip-components 1

RUN ./configure && make && make install

RUN git clone https://github.com/dloscutoff/ascii-piet.git

RUN chmod +x /ascii-piet/ascii2piet.py

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh /bin/touch   /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1      \
              /lib/libz.so.1                /lib/
COPY --from=0 /usr/local/bin/npiet          /usr/bin/
COPY --from=0 /ascii-piet/ascii2piet.py     /usr/bin/piet.py
COPY --from=0 /usr/bin/python3.13           /usr/bin/python3
COPY --from=0 /usr/lib/libX11.so.6          \
              /usr/lib/libXau.so.6          \
              /usr/lib/libXdmcp.so.6        \
              /usr/lib/libXpm.so.4          \
              /usr/lib/libaom.so.3          \
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
              /usr/lib/libxcb.so.1          \
              /usr/lib/libzstd.so.1         /usr/lib/
COPY --from=0 /usr/lib/python3.13           /usr/lib/python3.13

COPY /piet /usr/bin/

ENTRYPOINT ["piet"]
