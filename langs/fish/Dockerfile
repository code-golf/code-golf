FROM alpine:3.12 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base curl libffi-dev zlib-dev

RUN curl https://www.python.org/ftp/python/3.6.0/Python-3.6.0.tar.xz \
  | tar xJf -

RUN cd Python-3.6.0       \
 && ./configure           \
    --prefix=/usr         \
    --with-lto            \
 && make -j`proc` install \
 && strip /usr/bin/python3.6

ADD fish.py /

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/bin/python3.6       /usr/bin/python
COPY --from=0 /usr/lib/python3.6       /usr/lib/python3.6

COPY --from=0 /fish.py                 /usr/bin/fish

ENTRYPOINT ["/usr/bin/python", "./usr/bin/fish", "-c", "1n'.'o0n;"]
