FROM alpine:3.23 AS builder

RUN apk add --no-cache clang curl llvm-dev make openssl-dev

ENV PYVER=3.14.2 VER=70c9fd9

RUN curl -#L https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz \
  | tar xJ --strip-components 1

RUN CC='clang' ./configure     \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto=thin            \
    --with-tail-call-interp    \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.14  \
          /usr/lib/python3.14/lib-dynload/*.so

RUN curl -#L https://github.com/DennisMitchell/jellylanguage/tarball/$VER \
  | tar xz --directory /tmp --strip-components 1

# FIXME Python 3.14.X is denied.
RUN pip3 install --upgrade /tmp

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

COPY jelly.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/jellywrapper -s -static jelly.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1   /lib/
COPY --from=0 /usr/bin/jelly        \
              /usr/bin/jellywrapper \
              /usr/bin/python3.14   /usr/bin/
COPY --from=0 /usr/lib/libz.so.1    /usr/lib/
COPY --from=0 /usr/lib/python3.14   /usr/lib/python3.14
