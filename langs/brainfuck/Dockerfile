FROM alpine:3.13 as builder

RUN mkdir /empty

ENV PY_VERSION=3.7-v7.3.4 VERSION=61d1f49

RUN apk add --no-cache curl gcc libffi-dev linux-headers make musl-dev python2

RUN curl https://downloads.python.org/pypy/pypy$PY_VERSION-src.tar.bz2 \
  | tar xj

RUN curl -L https://github.com/primo-ppcg/bf-jit/tarball/$VERSION \
  | tar xz

RUN LDFLAGS=-static python pypy$PY_VERSION-src/rpython/bin/rpython \
    --lto --opt=jit primo-ppcg-bf-jit-$VERSION/bf-jit.py           \
 && strip bf-jit-c

FROM scratch

COPY --from=0 /empty    /proc
COPY --from=0 /empty    /tmp
COPY --from=0 /bf-jit-c /usr/bin/brainfuck
