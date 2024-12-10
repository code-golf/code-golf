FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl libffi-dev openssl-dev zlib-dev

# Python meta: https://pypi.org/project/vyxal/
ENV PYVER=3.10.15 VER=2.22.4.3

RUN curl https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

WORKDIR /Python-$PYVER

RUN ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.10  \
          /usr/lib/python3.10/lib-dynload/*.so

RUN pip3 install vyxal==$VER

WORKDIR /usr/lib/python3.10

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

COPY vyxalwrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/vyxalwrapper -s -static /vyxalwrapper.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/python3.10      \
              /usr/bin/vyxal           \
              /usr/bin/vyxalwrapper    /usr/bin/
COPY --from=0 /usr/lib/libffi.so.8     \
              /usr/lib/libz.so.1       /usr/lib/
COPY --from=0 /usr/lib/python3.10      /usr/lib/python3.10

ENTRYPOINT ["vyxalwrapper"]

CMD ["--version"]
