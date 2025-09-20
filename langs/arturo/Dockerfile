FROM nimlang/nim:2.2.4-alpine-slim AS builder

RUN apk add --no-cache pcre

ENV VER=51e8a5e

RUN curl -#L https://github.com/arturo-lang/arturo/tarball/$VER \
  | tar xz --strip-components 1

# HACK: See https://github.com/microsoft/mimalloc/issues/1065
RUN mkdir /usr/include/linux \
 && cp /usr/include/sys/prctl.h /usr/include/linux/

RUN ./build.nims --install --mode mini --release

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /root/.arturo/bin/arturo /usr/bin/
COPY --from=0 /usr/lib/libpcre.so.1    /usr/lib/

ENTRYPOINT ["arturo"]

CMD ["--version"]
