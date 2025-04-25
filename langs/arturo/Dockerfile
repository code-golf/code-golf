FROM nimlang/nim:2.2.2-alpine AS builder

RUN apk add --no-cache pcre

ENV VER=5f3443c

RUN curl -#L https://github.com/arturo-lang/arturo/tarball/$VER \
  | tar xz --strip-components 1

RUN ./build.nims --install --mode mini --release

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /root/.arturo/bin/arturo /usr/bin/
COPY --from=0 /usr/lib/libpcre.so.1    /usr/lib/

ENTRYPOINT ["arturo"]

CMD ["--version"]
