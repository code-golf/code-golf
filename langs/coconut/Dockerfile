FROM alpine:3.19 AS builder

RUN apk add --no-cache build-base curl libffi-dev zlib-dev

ENV PYVER 3.12.1
ENV VER 3.0.4

RUN curl https://www.python.org/ftp/python/$VER/Python-$PYVER.tar.xz | tar xJ

RUN cd Python-$PYVER           \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.12  \
 && pip install coconut==$VER

FROM codegolf/lang-base-no-proc

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/python3.12      /usr/bin/python
COPY --from=0 /usr/lib/python3.12      /usr/lib/python3.12
COPY --from=0 /usr/local/bin/coconut   /usr/local/bin/coconut

ENTRYPOINT ["coconut"]

CMD ["-v"]
