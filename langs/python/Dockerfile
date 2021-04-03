FROM alpine:3.13 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base curl libffi-dev zlib-dev

RUN curl https://www.python.org/ftp/python/3.9.3/Python-3.9.3.tar.xz \
  | tar xJf -

RUN cd Python-3.9.3       \
 && ./configure           \
    --prefix=/usr         \
    --with-lto            \
 && make -j`proc` install \
 && strip /usr/bin/python3.9

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/bin/python3.9       /usr/bin/python
COPY --from=0 /usr/lib/python3.9       /usr/lib/python3.9

ENTRYPOINT ["python", "-c", "import platform;print(platform.python_version())"]
