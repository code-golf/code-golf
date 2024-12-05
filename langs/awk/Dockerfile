FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl

RUN curl https://ftp.gnu.org/gnu/gawk/gawk-5.3.1.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure    \
    --disable-lint \
    --prefix /usr  \
 && make install   \
 && strip /usr/bin/gawk

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/gawk            /usr/bin/
COPY --from=0 /usr/lib/gawk            /usr/lib/gawk

ENTRYPOINT ["gawk"]

CMD ["--version"]
