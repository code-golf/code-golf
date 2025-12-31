FROM rust:1.92.0-alpine3.23 AS builder

RUN apk add --no-cache bash curl

ENV VER=0.5.1

WORKDIR /amber

RUN curl -#L https://github.com/amber-lang/amber/archive/refs/tags/$VER-alpha.tar.gz \
  | tar xz --strip-components 1

RUN cargo install \
    --path .      \
    --root /usr

FROM codegolf/lang-base

COPY --from=0 /bin/bash /bin/sed        /bin/
COPY --from=0 /lib/ld-musl-*.so.1       /lib/
COPY --from=0 /usr/bin/amber            \
              /usr/bin/bc               \
              /usr/bin/env              \
              /usr/bin/seq              /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6 \
              /usr/lib/libreadline.so.8 /usr/lib/

ENTRYPOINT ["amber"]

CMD ["--version"]
