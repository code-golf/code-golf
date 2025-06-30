FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl libffi-dev openssl-dev zlib-dev

ENV PYVER=3.13.5 VER=1.1.0

RUN curl -#L https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz \
  | tar xJ --strip-components 1

RUN CC='cc' ./configure        \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto=thin            \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.13  \
          /usr/lib/python3.13/lib-dynload/*.so

RUN pip3 install hy==$VER

WORKDIR /usr/lib/python3.13

# Remove some unneeded libraries.
RUN find . -name '*.opt-[12].pyc' -delete \
 && rm -r __pycache__/doctest.*           \
          __pycache__/pydoc.*             \
          __pycache__/turtle.*            \
          ensurepip                       \
          idlelib                         \
          pydoc_data                      \
          tkinter                         \
          turtle.py                       \
          turtledemo                      \
          unittest

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/hy              \
              /usr/bin/python3.13      /usr/bin/
COPY --from=0 /usr/lib/libz.so.1       /usr/lib/
COPY --from=0 /usr/lib/python3.13      /usr/lib/python3.13

ENTRYPOINT ["hy"]

CMD ["--version"]
