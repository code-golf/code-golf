FROM alpine:3.20 AS builder

ENV PYVER=3.12.6

RUN apk add --no-cache build-base curl git libffi-dev openssl-dev zlib-dev

RUN git clone https://github.com/DennisMitchell/jellylanguage.git

RUN curl https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

RUN cd Python-$PYVER           \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.12  \
          /usr/lib/python3.12/lib-dynload/*.so

RUN pip3 install --upgrade /jellylanguage/

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

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/
COPY --from=0 /usr/bin/jelly           \
              /usr/bin/python3.12      /usr/bin/
COPY --from=0 /usr/lib/python3.12      /usr/lib/python3.12

ENTRYPOINT ["jelly"]
