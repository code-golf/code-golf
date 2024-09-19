FROM crystallang/crystal:1.13.3-alpine AS builder

RUN apk add --no-cache build-base curl gc-dev libevent-dev libstdc++-dev llvm18-dev llvm18-static zstd-dev zstd-static

RUN curl -L https://github.com/crystal-lang/crystal/archive/refs/tags/1.13.3.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` crystal interpreter=1 release=1 static=1

RUN shards install

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /.build/crystal          /usr/bin/

COPY /crystalwrapper /usr/bin/

# Remove
COPY --from=0 /bin/chmod /bin/
RUN chmod +x /usr/bin/crystalwrapper

ENTRYPOINT ["crystalwrapper"]

CMD ["--version"]
