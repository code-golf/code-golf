FROM alpine:3.19 AS builder

RUN apk add --no-cache build-base curl libffi-dev zlib-dev openssl-dev

ENV PYVER 3.12.2
ENV VER 3.1.0

RUN curl https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

RUN cd Python-$PYVER           \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.12

RUN python3 -m pip install coconut==$VER

RUN rm -r /usr/lib/python3.12/site-packages/pip

FROM codegolf/lang-base-no-proc

COPY --from=0 /lib/ld-musl-*.so.1   \
              /lib/libc.musl-*.so.1 \
              /lib/libz.so.1        /lib/
COPY --from=0 /usr/bin/coconut      \
              /usr/bin/python3      \
              /usr/bin/python3.12   /usr/bin/
COPY --from=0 /usr/lib/python3.12   /usr/lib/python3.12/

ENTRYPOINT ["coconut"]

CMD ["-v"]
