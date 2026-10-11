FROM alpine:3.24 AS builder

RUN apk add --no-cache clang gcc llvm-dev make openssl-dev

ENV CC='clang' PYVER=3.15.0 VER=3.2.0

ADD --unpack https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz /

WORKDIR /Python-$PYVER

RUN ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto=thin            \
    --with-tail-call-interp    \
    --without-doc-strings      \
    --without-remote-debug     \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.15  \
          /usr/lib/python3.15/lib-dynload/*.so

RUN pip3 install --break-system-packages coconut==$VER

WORKDIR /usr/lib/python3.15

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
          turtledemo

FROM codegolf/lang-base

COPY --from=0 --parents /lib/ld-musl-*.so.1 \
                        /usr/bin/coconut    \
                        /usr/bin/python3.15 \
                        /usr/lib/libz.so*   \
                        /usr/lib/python3.15 /

ENTRYPOINT ["coconut"]

CMD ["--version"]
