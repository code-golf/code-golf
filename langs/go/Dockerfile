FROM golang:1.18.3-alpine3.16 as builder

RUN apk add --no-cache build-base gmp-dev

COPY go.c /go/

RUN gcc -s -o go go.c

FROM codegolf/lang-base

COPY --from=0 /go/go                                     /usr/bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                   /lib/
COPY --from=0 /usr/local/go/bin/go                       /usr/local/go/bin/
COPY --from=0 /usr/local/go/pkg/linux_amd64              /usr/local/go/pkg/linux_amd64/
COPY --from=0 /usr/local/go/pkg/tool/linux_amd64/compile \
              /usr/local/go/pkg/tool/linux_amd64/link    /usr/local/go/pkg/tool/linux_amd64/

ENTRYPOINT ["go"]

CMD ["version"]
