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

# Remove some unneeded libraries.
RUN cd /usr/lib/python3.13      \
 && rm -r __pycache__/doctest.* \
          __pycache__/pydoc.*   \
          __pycache__/turtle.*  \
          ensurepip             \
          idlelib               \
          pydoc_data            \
          tkinter               \
          turtle.py             \
          turtledemo

WORKDIR /osabie

RUN curl -#L https://github.com/Adriandmen/05AB1E/tarball/$VER \
  | tar xz --strip-components 1
        
COPY escape.patch /osabie/
        
RUN patch -p 1 < escape.patch

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh /bin/touch /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1    /lib/
COPY --from=0 /osabie/osabie.py           /osabie/
COPY --from=0 /osabie/behaviour           /osabie/behaviour
COPY --from=0 /osabie/lib                 /osabie/lib
COPY --from=0 /usr/bin/python3.13         /usr/bin/python
COPY --from=0 /usr/lib/python3.13         /usr/lib/python3.13

COPY 05ab1e /usr/bin/

ENTRYPOINT ["05ab1e"]
