FROM alpine:3.12 as builder

RUN apk add --no-cache curl gcc linux-headers make musl-dev python2

RUN curl https://downloads.python.org/pypy/pypy3.7-v7.3.2-src.tar.bz2 \
  | tar xj

RUN curl -L https://github.com/primo-ppcg/fish-jit/tarball/cee89b7 \
  | tar xz

RUN LDFLAGS=-static python pypy3.7-v7.3.2-src/rpython/bin/rpython \
    --lto primo-ppcg-fish-jit-cee89b7/fish-jit.py                 \
 && strip fish-jit-c

FROM scratch

COPY --from=0 fish-jit-c /usr/bin/fish

ENTRYPOINT ["/usr/bin/fish", "-c", "1n'.'o0n;"]
