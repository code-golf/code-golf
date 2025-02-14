FROM nimlang/nim:2.2.0-alpine-slim AS builder

RUN apk add --no-cache patch pcre

ENV VER=e57e8df

RUN curl -#L https://github.com/arturo-lang/arturo/tarball/$VER \
  | tar xz --strip-components 1

COPY stderr.patch /

RUN patch -p0 < stderr.patch

RUN ./build.nims \
    --install    \
    --mode mini  \
    --release

COPY arturo.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/arturo -s arturo.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/arturo          /usr/bin/
COPY --from=0 /usr/lib/libpcre.so.1    /usr/lib/
COPY --from=0 /root/.arturo/bin/arturo /usr/local/bin/

ENTRYPOINT ["arturo"]

CMD ["--version"]
