FROM golang:1.24.2-alpine3.21 AS builder

RUN apk add --no-cache gcc musl-dev

# Precompile the standard library.
RUN GOCACHE=/cache go install std

COPY gowrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/gowrapper -s -static /gowrapper.c

FROM codegolf/lang-base

COPY --from=0 /cache                   /cache
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/gowrapper       /usr/bin/
COPY --from=0 /usr/local/go            /usr/local/

ENTRYPOINT ["gowrapper"]

CMD ["--version"]
