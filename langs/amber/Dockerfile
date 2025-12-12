FROM rust:1.92.0-alpine3.23 AS builder

RUN apk add --no-cache bash curl

ENV CRT='x86_64-unknown-linux-musl' VER=0.5.1

WORKDIR /amber

RUN curl -#L https://github.com/amber-lang/amber/archive/refs/tags/$VER-alpha.tar.gz \
  | tar xz --strip-components 1

RUN rustup target add $CRT

RUN cargo install \
    --path .      \
    --root /usr   \
    --target $CRT

FROM codegolf/lang-base

# FIXME Target trixie and target.
#
#  > [builder 7/7] RUN ldd /usr/bin/amber && false:
# 0.449   statically linked
#
# ^ bullshit

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
