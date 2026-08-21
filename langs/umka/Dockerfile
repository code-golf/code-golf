FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev

ENV VER=1.5.6

WORKDIR /umka

ADD https://github.com/vtereshkov/umka-lang.git#v$VER .

RUN make -j`nproc` install

COPY umka.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/umka -s umka.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/umka       /usr/bin/
COPY --from=0 /usr/local/bin/umka /usr/local/bin/

ENTRYPOINT ["umka"]
