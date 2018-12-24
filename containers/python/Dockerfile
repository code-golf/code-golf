FROM alpine:edge as builder

RUN apk add --no-cache curl gcc libffi-dev linux-headers make musl-dev zlib-dev

RUN curl https://www.python.org/ftp/python/3.7.2/Python-3.7.2.tar.xz \
  | tar xJf -

RUN cd Python-3.7.2       \
 && ./configure           \
    --prefix=/usr         \
    --with-lto            \
 && make -j`proc` install \
 && strip /usr/bin/python3.7

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/python3.7       /usr/bin/python
COPY --from=0 /usr/lib/python3.7       /usr/lib/python3.7

ENTRYPOINT ["/usr/bin/python", "-c", "import platform;print(platform.python_version())"]
