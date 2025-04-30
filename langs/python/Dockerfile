FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl libffi-dev zlib-dev

# Also bump Coconut, Hy.
ENV VER=3.13.3

RUN curl -# https://www.python.org/ftp/python/$VER/Python-$VER.tar.xz | tar xJ

WORKDIR /Python-$VER

RUN ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-ensurepip        \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.13  \
          /usr/lib/python3.13/lib-dynload/*.so

WORKDIR /usr/lib/python3.13

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

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/python3.13      /usr/bin/python
COPY --from=0 /usr/lib/python3.13      /usr/lib/python3.13

ENTRYPOINT ["python"]

CMD ["--version"]
