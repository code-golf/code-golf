FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl libffi-dev openssl-dev zlib-dev

ENV PYVER=3.13.0 VER=3.1.2

RUN curl https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

RUN cd Python-$PYVER           \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.13  \
          /usr/lib/python3.13/lib-dynload/*.so

RUN pip3 install coconut==$VER

# Remove some unneeded libraries.
RUN rm -r /usr/lib/python3.13/__pycache__/doctest.* \
          /usr/lib/python3.13/__pycache__/pydoc.*   \
          /usr/lib/python3.13/__pycache__/turtle.*  \
          /usr/lib/python3.13/ensurepip             \
          /usr/lib/python3.13/idlelib               \
          /usr/lib/python3.13/pydoc_data            \
          /usr/lib/python3.13/tkinter               \
          /usr/lib/python3.13/turtle.py             \
          /usr/lib/python3.13/turtledemo

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/
COPY --from=0 /usr/bin/coconut         \
              /usr/bin/python3.13      /usr/bin/
COPY --from=0 /usr/lib/python3.13      /usr/lib/python3.13

ENTRYPOINT ["coconut"]

CMD ["-v"]
