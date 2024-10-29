FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl git libffi-dev openssl-dev py3-pillow zlib-dev

ENV VER=3.12.7

RUN curl https://www.python.org/ftp/python/$VER/Python-$VER.tar.xz | tar xJ

RUN cd Python-$VER             \
 && ./configure                \
    --disable-test-modules     \
    --prefix /usr              \
    --with-lto                 \
    --without-ensurepip        \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.12  \
          /usr/lib/python3.12/lib-dynload/*.so

# Remove some unneeded libraries.
RUN cd /usr/lib/python3.12      \
 && rm -r __pycache__/doctest.* \
          __pycache__/pydoc.*   \
          __pycache__/turtle.*  \
          ensurepip             \
          idlelib               \
          lib2to3               \
          pydoc_data            \
          tkinter               \
          turtle.py             \
          turtledemo

RUN git clone https://github.com/dloscutoff/ascii-piet.git

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh                /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1        \
              /lib/libz.so.1                  /lib/
COPY --from=0 /ascii-piet/ascii2piet.py       /usr/bin/
COPY --from=0 /usr/bin/python3.12             /usr/bin/python
COPY --from=0 /usr/lib/libXau.so.6            \
              /usr/lib/libXdmcp.so.6          \
              /usr/lib/libbsd.so.0            \
              /usr/lib/libgcc_s.so.1          \
              /usr/lib/libimagequant.so.0.0.4 \
              /usr/lib/libjpeg.so.8           \
              /usr/lib/libmd.so.0             \
              /usr/lib/libopenjp2.so.7        \
              /usr/lib/libsharpyuv.so.0       \
              /usr/lib/libtiff.so.6           \
              /usr/lib/libwebp.so.7           \
              /usr/lib/libxcb.so.1            \
              /usr/lib/libzstd.so.1           /usr/lib/
COPY --from=0 /usr/lib/python3.12             /usr/lib/python3.12

COPY /piet /usr/bin/

ENTRYPOINT ["piet"]
