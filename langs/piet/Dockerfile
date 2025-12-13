FROM alpine:3.23 AS builder

RUN apk add --no-cache clang curl gd-dev giflib-dev groff llvm-dev make openssl-dev zlib-dev

ENV CC='clang' PYVER=3.14.2 VER=1.3f

WORKDIR /python

RUN curl -#L https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto=thin            \
    --with-tail-call-interp    \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.14  \
          /usr/lib/python3.14/lib-dynload/*.so

RUN pip3 install pillow

WORKDIR /usr/lib/python3.14

# Remove some unneeded libraries.
RUN find . -name '*.opt-[12].pyc' -delete \
 && rm -r __pycache__/doctest.*           \
          __pycache__/pydoc.*             \
          __pycache__/turtle.*            \
          _pyrepl                         \
          ensurepip                       \
          idlelib                         \
          pydoc_data                      \
          tkinter                         \
          turtle.py                       \
          turtledemo                      \
          unittest

WORKDIR /piet

RUN curl -#L https://www.bertnase.de/npiet/npiet-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure && make && make install

ADD --chmod=777 https://raw.githubusercontent.com/dloscutoff/ascii-piet/refs/heads/master/ascii2piet.py /usr/bin/ascii-piet

COPY piet.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/piet -s -static /piet.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1           /lib/
COPY --from=0 /usr/bin/ascii-piet           \
              /usr/bin/piet                 /usr/bin/
COPY --from=0 /usr/bin/python3.14           /usr/bin/python3
COPY --from=0 /usr/lib/libX11.so.6          \
              /usr/lib/libXau.so.6          \
              /usr/lib/libXdmcp.so.6        \
              /usr/lib/libXpm.so.4          \
              /usr/lib/libaom.so.3          \
              /usr/lib/libavif.so.16        \
              /usr/lib/libbrotlicommon.so.1 \
              /usr/lib/libbrotlidec.so.1    \
              /usr/lib/libbsd.so.0          \
              /usr/lib/libbz2.so.1          \
              /usr/lib/libdav1d.so.7        \
              /usr/lib/libexpat.so.1        \
              /usr/lib/libfontconfig.so.1   \
              /usr/lib/libfreetype.so.6     \
              /usr/lib/libgcc_s.so.1        \
              /usr/lib/libgd.so.3           \
              /usr/lib/libgif.so.7          \
              /usr/lib/libjpeg.so.8         \
              /usr/lib/libmd.so.0           \
              /usr/lib/libpng16.so.16       \
              /usr/lib/libsharpyuv.so.0     \
              /usr/lib/libstdc++.so.6       \
              /usr/lib/libtiff.so.6         \
              /usr/lib/libwebp.so.7         \
              /usr/lib/libxcb.so.1          \
              /usr/lib/libyuv.so            \
              /usr/lib/libz.so.1            \
              /usr/lib/libzstd.so.1         /usr/lib/
COPY --from=0 /usr/lib/python3.14           /usr/lib/python3.14
COPY --from=0 /usr/local/bin/npiet          /usr/local/bin/

ENTRYPOINT ["piet"]
