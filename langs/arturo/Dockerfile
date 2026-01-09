FROM nimlang/nim:2.2.6-alpine-slim AS builder

ENV VER=c09869c

RUN curl -#L https://github.com/arturo-lang/arturo/tarball/$VER \
  | tar xz --strip-components 1

RUN ./build.nims --install --log --mode mini --release

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /bin/arturo         /usr/bin/

ENTRYPOINT ["arturo"]
