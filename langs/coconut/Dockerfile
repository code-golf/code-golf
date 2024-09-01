FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl libffi-dev openssl-dev zlib-dev

ENV PYVER=3.12.5 VER=3.1.1

RUN curl https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

RUN cd Python-$PYVER           \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.12

RUN pip3 install coconut==$VER

# Remove some unneeded libraries.
RUN cd /usr/lib/python3.12 \
 && rm -r ensurepip idlelib tkinter

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/
COPY --from=0 /usr/bin/coconut         \
              /usr/bin/python3.12      /usr/bin/
COPY --from=0 /usr/lib/python3.12      /usr/lib/python3.12

ENTRYPOINT ["coconut"]

CMD ["-v"]
