FROM golang:1.17.4-alpine3.15 as builder

RUN mkdir /empty

FROM scratch

COPY --from=0 /bin                                       /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                   /lib/
COPY --from=0 /empty                                     /proc
COPY --from=0 /empty                                     /tmp
COPY --from=0 /usr/local/go/bin/go                       /usr/local/go/bin/
COPY --from=0 /usr/local/go/pkg/linux_amd64              /usr/local/go/pkg/linux_amd64/
COPY --from=0 /usr/local/go/pkg/tool/linux_amd64/compile \
              /usr/local/go/pkg/tool/linux_amd64/link    /usr/local/go/pkg/tool/linux_amd64/

COPY go /usr/bin/go

ENTRYPOINT ["/usr/bin/go"]

CMD ["version"]
