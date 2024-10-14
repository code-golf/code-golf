FROM rust:1.81.0-alpine3.20 AS builder

RUN apk add --no-cache bash build-base curl

ENV VER=0.3.5

RUN curl -L https://github.com/amber-lang/amber/archive/refs/tags/$VER-alpha.tar.gz | tar xz

RUN cd amber-$VER-alpha \
 && cargo install --path .

FROM codegolf/lang-base

COPY --from=0 /bin/bash /bin/sed         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1   /lib/
COPY --from=0 /usr/bin/bc                \
              /usr/bin/env               \
              /usr/bin/seq               \
              /usr/local/cargo/bin/amber /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6  \
              /usr/lib/libreadline.so.8  /usr/lib/

ENTRYPOINT ["amber"]

CMD ["--version"]
