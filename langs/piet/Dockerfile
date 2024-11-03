FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl openssl-dev zlib-dev

ENV PYVER=3.13.0 VERSION=2ba1c90

RUN curl https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

RUN cd Python-$PYVER           \
 && ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-openssl=/usr        \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.13

RUN pip3 install Pillow

RUN curl -L https://github.com/dloscutoff/ascii-piet/tarball/$VERSION \
  | tar xz --strip-components 1

COPY /piet.c /

RUN gcc -Wall -Werror -Wextra -s -o /piet /piet.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/
COPY --from=0 /ascii2piet.py /piet     /usr/bin/
COPY --from=0 /usr/bin/python3.13      /usr/bin/python
COPY --from=0 /usr/lib/python3.13      /usr/lib/python3.13

ENTRYPOINT ["piet"]
