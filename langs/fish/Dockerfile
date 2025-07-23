FROM alpine:3.22 AS builder

RUN apk add -X https://dl-cdn.alpinelinux.org/alpine/edge/testing \
    --no-cache curl gcc libffi-dev linux-headers make musl-dev pypy

ENV LDFLAGS='-static' PYVER=3.11-v7.3.19 VER=96b1cfc

RUN curl -#L https://downloads.python.org/pypy/pypy$PYVER-src.tar.bz2 \
  | tar xj --strip-components 1

WORKDIR /usr/bin

RUN curl -#L https://github.com/primo-ppcg/fish-jit/tarball/$VER \
  | tar xz --strip-components 1

RUN pypy /rpython fish-jit \
    --lto                  \
    --make-jobs=`nproc`    \
    --opt=jit              \
 && strip fish-jit-c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/fish-jit-c /usr/bin/fish

ENTRYPOINT ["fish"]

CMD ["--no-prng", "-u"]
