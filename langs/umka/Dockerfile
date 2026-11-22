FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev

ENV VER=1.5.7

WORKDIR /umka

ADD https://github.com/vtereshkov/umka-lang.git#v$VER .

RUN make -j`nproc` install PREFIX=/usr

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/umka       /usr/bin/

ENTRYPOINT ["umka"]
