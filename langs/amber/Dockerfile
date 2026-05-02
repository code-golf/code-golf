FROM rust:1.95.0-alpine3.23 AS builder

RUN apk add --no-cache bash

ENV VER=0.6.0

WORKDIR /amber

ADD https://github.com/amber-lang/amber.git#$VER-alpha .

RUN cargo install \
    --path .      \
    --root /usr/local

COPY amber.c .

RUN gcc -Wall -Werror -Wextra -o /usr/bin/amber -s amber.c

FROM codegolf/lang-base

COPY --from=0 /bin/bash                 /bin/
COPY --from=0 /lib/ld-musl-*.so.1       /lib/
COPY --from=0 /usr/bin/amber            \
              /usr/bin/env              /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6 \
              /usr/lib/libreadline.so.8 /usr/lib/
COPY --from=0 /usr/local/bin/amber      /usr/local/bin/

ENTRYPOINT ["amber"]

CMD ["--version"]
