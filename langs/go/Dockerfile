ARG VER=1.26.3

FROM golang:$VER-alpine3.23 AS builder

RUN apk add --no-cache gcc musl-dev

# Precompile the standard library.
RUN GOCACHE='/cache' go install std

COPY go.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/gowrapper -s -static /go.c

FROM codegolf/lang-base

COPY --from=0 --parents /cache                                     \
                        /lib/ld-musl-*.so.1                        \
                        /usr/bin/gowrapper                         \
                        /usr/local/go/bin/go                       \
                        /usr/local/go/pkg/tool/linux_amd64/asm     \
                        /usr/local/go/pkg/tool/linux_amd64/compile \
                        /usr/local/go/pkg/tool/linux_amd64/link    \
                        /usr/local/go/src                          /

ENTRYPOINT ["gowrapper"]

CMD ["--version"]
