FROM rust:1.84.1-alpine3.21 AS builder

RUN apk add --no-cache build-base linux-headers

RUN cargo install uiua --root /usr --version 0.14.1 \
 && strip /usr/bin/uiua

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/uiua            /usr/bin/

ENTRYPOINT ["uiua"]

CMD ["--version"]
