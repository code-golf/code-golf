FROM rust:1.80.1-alpine3.20 AS builder

RUN apk add --no-cache bash build-base curl

RUN curl -L https://github.com/amber-lang/amber/archive/refs/tags/0.3.5-alpha.tar.gz \
  | tar xz --strip-components 1

RUN cargo install --path .

FROM codegolf/lang-base

COPY --from=0 /bin/bash                  \
              /bin/busybox               \
              /bin/sed                   /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1   /lib/
COPY --from=0 /usr/bin/bc                \
              /usr/bin/env               \
              /usr/bin/seq               \
              /usr/local/cargo/bin/amber /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6  \
              /usr/lib/libreadline.so.8  /usr/lib/

ENTRYPOINT ["amber"]

CMD ["--version"]
