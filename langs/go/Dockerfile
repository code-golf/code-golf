FROM golang:1.14.2-alpine as builder

FROM scratch

COPY --from=0 / /

COPY run-go /usr/bin/run-go

ENTRYPOINT ["/usr/local/go/bin/go", "version"]
