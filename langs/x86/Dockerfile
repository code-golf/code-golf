FROM alpine:3.10 as builder

RUN apk add --no-cache curl gcc make musl-dev

RUN curl -L https://www.nasm.us/pub/nasm/releasebuilds/2.14.02/nasm-2.14.02.tar.gz | tar xzf -

RUN mv nasm-* nasm   \
    && cd nasm       \
    && ./configure   \
       --prefix=/usr \
    && make          \
    && make install  \
    && strip /usr/bin/nasm

COPY x86.c /

RUN gcc -s -o x86 x86.c

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /x86 /usr/bin/nasm       /usr/bin/
COPY header.asm                        /

ENTRYPOINT ["/usr/bin/nasm", "-v"]
