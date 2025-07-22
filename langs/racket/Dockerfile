FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl linux-headers

ENV VER=8.17

RUN curl -#L https://github.com/racket/racket/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN make bc -j`nproc` PKGS='' RACKETBC_SUFFIX=''

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /racket/bin/racket       /usr/bin/
COPY --from=0 /racket/collects         /usr/collects

ENTRYPOINT ["racket"]

CMD ["--version"]
