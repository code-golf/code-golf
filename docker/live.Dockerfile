FROM golang:1.23.3-alpine3.20

ENV CGO_ENABLED=0 GOAMD64=v4 GOPATH=

RUN apk add --no-cache brotli build-base linux-headers npm tzdata zopfli

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN ./esbuild \
 && find dist \( -name '*.css' -or  -name '*.js' -or -name '*.map' \) \
  | xargs -i -n1 -P`nproc` sh -c 'brotli {} && zopfli {}'

RUN go build -ldflags -s -trimpath \
 && gcc -Wall -Werror -Wextra -o /usr/bin/run-lang -s -static run-lang.c

FROM scratch

COPY --from=codegolf/lang-stax         ["/", "/langs/stax/rootfs/"        ] #    0 MiB

COPY --from=0 /go/code-golf                      /
COPY --from=0 /go/dist                           /dist/
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY          /public                            /public/
COPY --from=0 /usr/bin/run-lang                  /usr/bin/
COPY --from=0 /usr/share/zoneinfo                /usr/share/zoneinfo/

CMD ["/code-golf"]
