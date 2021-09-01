FROM alpine:3.14 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base curl libffi-dev zlib-dev

RUN curl https://www.python.org/ftp/python/3.9.6/Python-3.9.6.tar.xz | tar xJ

RUN cd Python-3.9.6       \
 && ./configure           \
    --prefix=/usr         \
    --with-lto            \
    --without-ensurepip   \
 && make -j`proc` install \
 && strip /usr/bin/python3.9

# Remove tests and .a files, but keep .pyc files to aid execution speed.
# https://github.com/docker-library/python/blob/master/Dockerfile-alpine.template
RUN find /usr/lib/python3.9 -depth \( \
        \( -type d -a \( -name idle_test -o -name test -o -name tests \) \) \
     -o \( -type f -a \( -name '*.a' \) \) \
    \) -exec rm -rf {} +

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /empty                   /proc
COPY --from=0 /empty                   /tmp
COPY --from=0 /usr/bin/python3.9       /usr/bin/python
COPY --from=0 /usr/lib/python3.9       /usr/lib/python3.9

ENTRYPOINT ["python"]

CMD ["-c", "import platform;print(platform.python_version())"]
