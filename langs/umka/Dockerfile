FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV VER=1.5.3

RUN curl -#L https://github.com/vtereshkov/umka-lang/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` install

COPY umka.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/umka -s umka.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/umka            /usr/bin/
COPY --from=0 /usr/local/bin/umka      /usr/local/bin/

ENTRYPOINT ["umka"]
