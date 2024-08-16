FROM alpine:3.20 AS builder

RUN apk add --no-cache cargo

RUN cargo install hush --root /usr --version 0.1.4

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/hush            /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1   /usr/lib/

ENTRYPOINT ["hush"]

CMD ["--version"]
