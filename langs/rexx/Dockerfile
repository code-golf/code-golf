FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev

ENV VER=3.9.7

ADD --unpack https://downloads.sourceforge.net/project/regina-rexx/regina-rexx/$VER/regina-rexx-$VER.tar.gz /

WORKDIR /regina-rexx-$VER

RUN ./configure             \
    --disable-flists        \
    --disable-gci           \
    --disable-posix-threads \
    --disable-screen        \
    --enable-64bit          \
    --prefix=/usr           \
 && make install

FROM codegolf/lang-base

COPY --from=0 --parents /lib/ld-musl-*.so.1 /usr/bin/rexx /

ENTRYPOINT ["rexx"]

CMD ["--version"]
