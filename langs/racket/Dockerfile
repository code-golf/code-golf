FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc linux-headers make musl-dev

ENV VER=9.3

WORKDIR /racket

ADD https://github.com/racket/racket.git#v$VER .

RUN make bc -j`nproc` PKGS='' RACKETBC_SUFFIX=''

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1       /lib/
COPY --from=0 /racket/racket/bin/racket /usr/bin/
COPY --from=0 /racket/racket/collects   /usr/collects

ENTRYPOINT ["racket"]

CMD ["--version"]
