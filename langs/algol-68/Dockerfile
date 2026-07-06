FROM alpine:3.24 AS builder

RUN apk add --no-cache clang gcc make patch

ENV LDFLAGS='-static' VER=3.12.3

ADD --unpack https://algol68genie.nl/algol68g-$VER.tar.gz /

WORKDIR /algol68g-$VER

COPY algol-68.patch .

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
