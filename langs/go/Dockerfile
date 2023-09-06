FROM golang:1.21.1-alpine3.18 AS builder

# Precompile the standard library.
RUN GOCACHE=/cache go install std

FROM codegolf/lang-base

COPY --from=0 /bin                     /bin
COPY --from=0 /cache                   /cache
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/local/go            /usr/local/go

COPY go /usr/bin/

ENTRYPOINT ["go"]

CMD ["version"]
