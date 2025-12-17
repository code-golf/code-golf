ARG VER=1.25.5

FROM golang:$VER-alpine3.23 AS builder

RUN apk add --no-cache gcc musl-dev

ENV GOCACHE='/cache'

# Precompile the standard library.
RUN go install std

COPY go.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/go -s -static go.c

FROM codegolf/lang-base

COPY --chown=nobody                                      \
     --from=0 /cache                                     /cache
COPY --from=0 /lib/ld-musl-*.so.1                        /lib/
COPY --from=0 /usr/bin/go                                /usr/bin/
COPY --from=0 /usr/local/go/bin/go                       /usr/local/go/bin/
COPY --from=0 /usr/local/go/pkg/tool/linux_amd64/asm     \
              /usr/local/go/pkg/tool/linux_amd64/compile \
              /usr/local/go/pkg/tool/linux_amd64/link    /usr/local/go/pkg/tool/linux_amd64/
COPY --from=0 /usr/local/go/src                          /usr/local/go/src

ENTRYPOINT ["go"]

CMD ["--version"]
