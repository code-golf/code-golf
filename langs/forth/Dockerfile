FROM alpine:3.18 AS builder

RUN apk add --no-cache build-base curl m4

RUN curl https://ftp.gnu.org/gnu/gforth/gforth-0.7.3.tar.gz | tar xz

RUN cd gforth-* \
 && ./configure --prefix=/usr --without-check \
 && make \
 && make install

FROM codegolf/lang-base

COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=builder /usr/bin/gforth          /usr/bin/forth
COPY --from=builder /usr/include/gforth      /usr/include/gforth
COPY --from=builder /usr/lib/gforth          /usr/lib/gforth
COPY --from=builder /usr/share/gforth        /usr/share/gforth

ENTRYPOINT ["forth"]

CMD ["--version"]
