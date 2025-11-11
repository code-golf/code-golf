FROM alpine:3.22 AS builder

RUN apk add --no-cache clang curl libffi-dev llvm make zlib-dev

# Also bump Coconut, Hy.
ENV VER=3.14.0

RUN curl -#L https://www.python.org/ftp/python/$VER/Python-$VER.tar.xz \
  | tar xJ --strip-components 1

RUN CC=clang ./configure       \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto=thin            \
    --with-tail-call-interp    \
    --without-ensurepip        \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.14  \
          /usr/lib/python3.14/lib-dynload/*.so

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

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/python3.14 /usr/bin/python
COPY --from=0 /usr/lib/python3.14 /usr/lib/python3.14

ENTRYPOINT ["python"]

CMD ["--version"]
