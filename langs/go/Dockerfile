FROM golang:1.23.4-alpine3.20 AS builder

# Precompile the standard library.
RUN GOCACHE=/cache go install std

FROM codegolf/lang-base

COPY --from=0 /bin                     /bin
COPY --from=0 /cache                   /cache
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/local/go            /usr/local/

COPY /gowrapper /usr/bin/

ENTRYPOINT ["gowrapper"]

CMD ["--version"]
