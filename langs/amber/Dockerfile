FROM rust:1.95.0-alpine3.23 AS builder

RUN apk add --no-cache bash

ENV VER=0.6.0

WORKDIR /amber

ADD https://github.com/amber-lang/amber.git#$VER-alpha .

RUN cargo install \
    --path .      \
    --root /usr

FROM codegolf/lang-base

COPY --from=0 /bin/bash                 /bin/
COPY --from=0 /lib/ld-musl-*.so.1       /lib/
COPY --from=0 /usr/bin/amber            /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6 \
              /usr/lib/libreadline.so.8 /usr/lib/

ENTRYPOINT ["amber"]

CMD ["--version"]
