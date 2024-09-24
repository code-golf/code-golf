FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl libffi-dev openssl-dev zlib-dev

ENV PYVER=3.13.2 VER=70c9fd9

RUN curl -# https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

WORKDIR /Python-$PYVER

RUN ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.13  \
          /usr/lib/python3.13/lib-dynload/*.so

WORKDIR /jelly

RUN curl -#L https://github.com/DennisMitchell/jellylanguage/tarball/$VER \
  | tar xz --strip-components 1

RUN pip3 install --upgrade .

WORKDIR /usr/lib/python3.13

# Remove some unneeded libraries.
RUN find . -name '*.opt-[12].pyc' -delete \
 && rm -r __pycache__/doctest.*           \
          __pycache__/pydoc.*             \
          __pycache__/turtle.*            \
          _pyrepl                         \
          ensurepip                       \
          idlelib                         \
          pydoc_data                      \
          tkinter                         \
          turtle.py                       \
          turtledemo                      \
          unittest

COPY jellywrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/jellywrapper -s -static /jellywrapper.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/jelly           \
              /usr/bin/jellywrapper    \
              /usr/bin/python3.13      /usr/bin/
COPY --from=0 /usr/lib/libz.so.1       /usr/lib/
COPY --from=0 /usr/lib/python3.13      /usr/lib/python3.13

ENTRYPOINT ["jellywrapper"]

CMD ["--version"]
