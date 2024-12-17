FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl libffi-dev openssl-dev zlib-dev

ENV PYVER=3.13.0 VERSION=70c9fd9

RUN curl -L https://github.com/DennisMitchell/jellylanguage/tarball/$VERSION | tar xz

RUN curl https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

RUN cd Python-$PYVER           \
 && ./configure                \
    --disable-test-modules     \
    --prefix /usr              \
    --with-lto                 \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.13  \
          /usr/lib/python3.13/lib-dynload/*.so

# Remove some unneeded libraries.
RUN cd /usr/lib/python3.13      \
 && rm -r __pycache__/doctest.* \
          __pycache__/pydoc.*   \
          __pycache__/turtle.*  \
          ensurepip             \
          idlelib               \
          pydoc_data            \
          tkinter               \
          turtle.py             \
          turtledemo

RUN pip3 install --upgrade /DennisMitchell-jellylanguage-$VERSION/

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sed /bin/sh /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1  \
              /lib/libz.so.1            /lib/
COPY --from=0 /usr/bin/jelly            \
              /usr/bin/python3.13       /usr/bin/
COPY --from=0 /usr/lib/python3.13       /usr/lib/python3.13

COPY /jellywrapper /usr/bin/

ENTRYPOINT ["jellywrapper"]
