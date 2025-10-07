ARG VER=1.25.1

FROM golang:$VER-alpine3.22 AS builder

RUN apk add --no-cache gcc musl-dev

# Precompile the standard library.
RUN GOCACHE=/cache go install std

COPY go.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/gowrapper -s -static /go.c

FROM codegolf/lang-base

COPY --from=0 /cache                                     /cache
COPY --from=0 /lib/ld-musl-x86_64.so.1                   /lib/
COPY --from=0 /usr/bin/gowrapper                         /usr/bin/
COPY --from=0 /usr/local/go/bin/go                       /usr/local/bin/
COPY --from=0 /usr/local/go/pkg/tool/linux_amd64/asm     \
              /usr/local/go/pkg/tool/linux_amd64/compile \
              /usr/local/go/pkg/tool/linux_amd64/link    /usr/local/pkg/tool/linux_amd64/
COPY --from=0 /usr/local/go/src                          /usr/local/src

ENTRYPOINT ["gowrapper"]

CMD ["--version"]
