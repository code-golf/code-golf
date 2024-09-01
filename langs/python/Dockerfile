FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl libffi-dev zlib-dev

# Also bump Coconut, Hy.
ENV VER=3.12.5

RUN curl https://www.python.org/ftp/python/$VER/Python-$VER.tar.xz | tar xJ

RUN cd Python-$VER             \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-ensurepip        \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.12  \
          /usr/lib/python3.12/lib-dynload/*.so

# Remove some unneeded libraries.
RUN rm -r /usr/lib/python3.12/__pycache__/doctest.* \
          /usr/lib/python3.12/__pycache__/pydoc.*   \
          /usr/lib/python3.12/__pycache__/turtle.*  \
          /usr/lib/python3.12/ensurepip             \
          /usr/lib/python3.12/idlelib               \
          /usr/lib/python3.12/lib2to3               \
          /usr/lib/python3.12/pydoc_data            \
          /usr/lib/python3.12/tkinter               \
          /usr/lib/python3.12/turtle.py             \
          /usr/lib/python3.12/turtledemo            \
          /usr/lib/python3.12/unittest

FROM codegolf/lang-base-no-proc

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/python3.12      /usr/bin/python
COPY --from=0 /usr/lib/python3.12      /usr/lib/python3.12

ENTRYPOINT ["python"]

CMD ["--version"]
