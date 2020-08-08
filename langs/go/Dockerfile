FROM golang:1.14.7-alpine3.12 as builder

FROM scratch

COPY --from=0 /bin                                       /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                   /lib/
COPY --from=0 /tmp                                       /tmp/
COPY --from=0 /usr/local/go/bin/go                       /usr/local/go/bin/
COPY --from=0 /usr/local/go/pkg/linux_amd64              /usr/local/go/pkg/linux_amd64/
COPY --from=0 /usr/local/go/pkg/tool/linux_amd64/compile \
              /usr/local/go/pkg/tool/linux_amd64/link    /usr/local/go/pkg/tool/linux_amd64/

COPY go /usr/bin/go

ENTRYPOINT ["/usr/local/go/bin/go", "version"]
