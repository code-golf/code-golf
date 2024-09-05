FROM crystallang/crystal:1.13.2-alpine AS builder

RUN apk add --no-cache curl libstdc++-dev llvm18-dev

RUN curl -L https://github.com/crystal-lang/crystal/archive/refs/tags/1.13.2.tar.gz \
  | tar xz --strip-components 1

RUN make crystal -j`nproc` release=1 interpreter=1 static=1

FROM codegolf/lang-base

COPY --from=0 /usr/bin/crystal /usr/bin/

ENTRYPOINT ["crystal"]

CMD ["-v"]
