FROM alpine:3.20 AS builder

ENV PY_VERSION=3.10-v7.3.17 VERSION=96b1cfc

RUN echo @edge https://dl-cdn.alpinelinux.org/alpine/edge/testing \
 >> /etc/apk/repositories                                         \
 && apk add --no-cache curl gcc libffi-dev linux-headers make musl-dev pypy@edge

RUN curl https://downloads.python.org/pypy/pypy$PY_VERSION-src.tar.bz2 \
  | tar xj

RUN curl -L https://github.com/primo-ppcg/fish-jit/tarball/$VERSION \
  | tar xz

RUN LDFLAGS=-static pypy pypy$PY_VERSION-src/rpython/bin/rpython \
    --lto --opt=jit primo-ppcg-fish-jit-$VERSION/fish-jit.py     \
 && strip fish-jit-c

FROM codegolf/lang-base

COPY --from=0 /fish-jit-c /usr/bin/fish

ENTRYPOINT ["fish"]

CMD ["--no-prng", "-u"]
