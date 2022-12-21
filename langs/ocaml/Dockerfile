FROM alpine:3.17 as builder

RUN apk add --no-cache ocaml=4.14.0-r0

FROM codegolf/lang-base

COPY --from=0 /usr/bin /usr/bin
COPY --from=0 /usr/lib/ocaml /usr/lib/ocaml
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/

ENTRYPOINT ["/usr/bin/ocaml"]

CMD ["--version"]
