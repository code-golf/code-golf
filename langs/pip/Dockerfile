FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl libffi-dev zlib-dev

ENV PYVER=3.12.5 VER=1.2.0

RUN curl https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

RUN cd Python-$PYVER           \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-ensurepip        \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.12

WORKDIR /tmp/dir

RUN curl -L https://github.com/dloscutoff/pip/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

FROM codegolf/lang-base

COPY --from=0 /bin/chmod /bin/sh       /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /tmp/dir                 /usr/bin
COPY --from=0 /usr/bin/python3.12      /usr/bin/python
COPY --from=0 /usr/lib/python3.12      /usr/lib/python3.12

RUN chmod +x /usr/bin/pip.py

ENTRYPOINT ["pip.py"]

CMD ["--version"]
