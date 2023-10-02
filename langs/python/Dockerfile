FROM alpine:3.18 AS builder

RUN apk add --no-cache build-base curl libffi-dev zlib-dev

ENV VER 3.12.0

RUN curl https://www.python.org/ftp/python/$VER/Python-$VER.tar.xz | tar xJ

RUN cd Python-$VER             \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-ensurepip        \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.12

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/python3.12      /usr/bin/python
COPY --from=0 /usr/lib/python3.12      /usr/lib/python3.12

ENTRYPOINT ["python"]

CMD ["-c", "import platform;print(platform.python_version())"]
