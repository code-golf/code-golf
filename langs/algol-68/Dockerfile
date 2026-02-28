FROM alpine:3.23 AS builder

RUN apk add --no-cache clang curl make patch

ENV LDFLAGS='-static' VER=3.10.10

RUN curl -#L https://jmvdveer.home.xs4all.nl/algol68g-$VER.tar.gz \
  | tar xz --strip-components 1

COPY algol-68.patch /

RUN patch -p0 < algol-68.patch

RUN ./configure             \
    --enable-arch=x86-64-v3 \
    --enable-clang          \
    --enable-generic        \
    --enable-stable         \
    --prefix=/usr           \
 && make install            \
 && strip /usr/bin/a68g

FROM codegolf/lang-base

COPY --from=0 /usr/bin/a68g /usr/bin/algol-68

ENTRYPOINT ["algol-68"]

CMD ["--version"]
