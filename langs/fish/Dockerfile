FROM alpine:3.24 AS builder

RUN apk add -X https://dl-cdn.alpinelinux.org/alpine/edge/testing \
    --no-cache gcc libffi-dev linux-headers make musl-dev pypy

ENV LDFLAGS='-static' PYVER=3.12-v8.0.0 VER=e9fe296

ADD --unpack https://downloads.python.org/pypy/pypy$PYVER-src.tar.gz /

WORKDIR /fish

ADD https://github.com/primo-ppcg/fish-jit.git?commit=$VER .

RUN pypy /pypy$PYVER-src/rpython fish-jit \
    --lto                                 \
    --make-jobs=`nproc`                   \
    --opt=jit                             \
 && strip fish-jit-c

FROM codegolf/lang-base

COPY --from=0 /fish/fish-jit-c /usr/bin/fish

ENTRYPOINT ["fish"]

CMD ["--no-prng", "-u"]
