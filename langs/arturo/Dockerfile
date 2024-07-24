FROM nimlang/nim:2.0.8-alpine-slim AS builder

RUN apk add --no-cache build-base pcre

RUN curl -L https://github.com/arturo-lang/arturo/archive/refs/tags/v0.9.83.tar.gz | tar xz

COPY build.patch /arturo-0.9.83/

WORKDIR /arturo-0.9.83

RUN patch -p 1 < build.patch

RUN ./build.nims install mini log

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  \
              /lib/libcrypto.so.3       \
              /lib/libssl.so.3          /lib/
COPY --from=0 /root/.arturo/bin/arturo  /usr/bin/
COPY --from=0 /usr/lib/libpcre.so.1     /usr/lib/

ENTRYPOINT ["arturo"]

CMD ["-v"]
