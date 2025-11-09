FROM nimlang/nim:2.2.6-alpine-slim AS builder

RUN apk add --no-cache pcre-dev

ENV VER=0991c79

RUN curl -#L https://github.com/arturo-lang/arturo/tarball/$VER \
  | tar xz --strip-components 1

RUN ./build.nims --install --log --mode mini --release

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1   /lib/
COPY --from=0 /bin/arturo           /usr/bin/
COPY --from=0 /usr/lib/libpcre.so.1 /usr/lib/

ENTRYPOINT ["arturo"]
