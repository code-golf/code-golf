# FIXME Staying on 3.15 for python2.
FROM alpine:3.15 AS builder

ENV PY_VERSION=3.9-v7.3.11 VERSION=888b90f

RUN apk add --no-cache curl gcc libffi-dev linux-headers make musl-dev python2

RUN curl https://downloads.python.org/pypy/pypy$PY_VERSION-src.tar.bz2 \
  | tar xj

RUN curl -L https://github.com/primo-ppcg/fish-jit/tarball/$VERSION \
  | tar xz

RUN LDFLAGS=-static python pypy$PY_VERSION-src/rpython/bin/rpython \
    --lto --opt=jit primo-ppcg-fish-jit-$VERSION/fish-jit.py       \
 && strip fish-jit-c

FROM codegolf/lang-base

COPY --from=0 /fish-jit-c /usr/bin/fish

ENTRYPOINT ["fish"]

CMD ["--no-prng", "-u"]
