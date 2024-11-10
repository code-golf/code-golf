FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl libffi-dev zlib-dev

# Also bump Coconut, Hy, Piet.
ENV VER=3.13.0

RUN curl https://www.python.org/ftp/python/$VER/Python-$VER.tar.xz | tar xJ

RUN cd Python-$VER             \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-ensurepip        \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.13  \
          /usr/lib/python3.13/lib-dynload/*.so

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

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/python3.13      /usr/bin/python
COPY --from=0 /usr/lib/python3.13      /usr/lib/python3.13

ENTRYPOINT ["python"]

CMD ["--version"]
