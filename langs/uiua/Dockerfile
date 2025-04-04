FROM rust:1.86.0-alpine3.21 AS builder

RUN apk add --no-cache build-base linux-headers

# 0.15.0 + patch broken zip2 dependency
RUN cargo install uiua --root /usr --git https://github.com/uiua-lang/uiua --rev f177aeb \
 && strip /usr/bin/uiua

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/uiua            /usr/bin/

ENTRYPOINT ["uiua"]

CMD ["--version"]
