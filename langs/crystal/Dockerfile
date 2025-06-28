ARG VER=1.16.3

FROM crystallang/crystal:$VER-alpine AS builder

ARG VER

RUN apk add --no-cache build-base curl gc-dev libevent-dev libstdc++-dev llvm18-dev llvm18-static zstd-dev zstd-static

RUN curl -#L https://github.com/crystal-lang/crystal/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` crystal interpreter=1 release=1 static=1

RUN strip /.build/crystal \
 && shards install

FROM codegolf/lang-base

COPY --from=0 /.build/crystal /usr/bin/

ENTRYPOINT ["crystal"]

CMD ["--version"]
