FROM alpine:3.24 AS builder

RUN apk add --no-cache clang gcc llvm-dev make

# Also bump Coconut, Hy.
ENV CC=clang VER=3.15.2

ADD --unpack https://www.python.org/ftp/python/$VER/Python-$VER.tar.xz /

WORKDIR /Python-$VER

RUN ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto=thin            \
    --with-tail-call-interp    \
    --without-doc-strings      \
    --without-ensurepip        \
    --without-remote-debug     \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.15  \
          /usr/lib/python3.15/lib-dynload/*.so

WORKDIR /usr/lib/python3.15

# Remove some unneeded libraries.
RUN find . -name '*.opt-[12].pyc' -delete   \
 && rm -r __pycache__/doctest.*             \
          __pycache__/pydoc.*               \
          __pycache__/turtle.*              \
          _pyrepl                           \
          ensurepip                         \
          idlelib                           \
          lib-dynload/binascii.cpython-*.so \
          lib-dynload/zlib.cpython-*.so     \
          profiling                         \
          pydoc_data                        \
          tkinter                           \
          turtle.py                         \
          turtledemo                        \
          unittest

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/python3.15 /usr/bin/python
COPY --from=0 /usr/lib/python3.15 /usr/lib/python3.15

ENTRYPOINT ["python"]

CMD ["--version"]
