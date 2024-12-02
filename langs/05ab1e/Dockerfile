FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl libffi-dev zlib-dev

ENV PYVER=3.13.0 VER=fe1ab18

RUN curl https://www.python.org/ftp/python/$PYVER/Python-$PYVER.tar.xz | tar xJ

WORKDIR /Python-$PYVER

RUN ./configure                \
    --disable-test-modules     \
    --prefix=/usr              \
    --with-lto                 \
    --without-ensurepip        \
    --without-static-libpython \
 && make -j`nproc` install     \
 && strip /usr/bin/python3.13  \
          /usr/lib/python3.13/lib-dynload/*.so

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

WORKDIR /osabie

RUN curl -#L https://github.com/Adriandmen/05AB1E/tarball/$VER \
  | tar xz --strip-components 1
        
COPY 05ab1e.c escape.patch /
        
RUN patch -p 1 < /escape.patch \
 && gcc -Wall -Werror -Wextra -o /usr/bin/05ab1e -s -static /05ab1e.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /osabie/osabie.py        /osabie/
COPY --from=0 /osabie/behaviour        /osabie/behaviour
COPY --from=0 /osabie/lib              /osabie/lib
COPY --from=0 /usr/bin/05ab1e          /usr/bin/
COPY --from=0 /usr/bin/python3.13      /usr/bin/osabie
COPY --from=0 /usr/lib/python3.13      /usr/lib/python3.13

ENTRYPOINT ["05ab1e"]

CMD ["--version"]
