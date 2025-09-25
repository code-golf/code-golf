FROM nimlang/nim:1.6.18-alpine-slim AS builder

RUN apk add --no-cache patch pcre

ENV VER=0.9.83

RUN curl -#L https://github.com/arturo-lang/arturo/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

COPY errors.patch /

RUN patch -p0 < errors.patch

RUN ./build.nims install mini release

COPY arturo.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/arturo -s -static arturo.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/arturo          /usr/bin/
COPY --from=0 /usr/lib/libcrypto.so.3  \
              /usr/lib/libpcre.so.1    \
              /usr/lib/libssl.so.3     /usr/lib/
COPY --from=0 /root/.arturo/bin/arturo /usr/local/bin/

ENTRYPOINT ["arturo"]

CMD ["--version"]
