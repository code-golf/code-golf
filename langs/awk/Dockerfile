FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl

RUN curl https://ftp.gnu.org/gnu/gawk/gawk-5.3.0.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure    \
    --disable-docs \
    --disable-lint \
 && make -j`nproc` \
 && make install   \
 && strip /usr/local/bin/gawk

FROM codegolf/lang-base-no-proc

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/local/bin/gawk      /usr/bin/

ENTRYPOINT ["gawk"]

CMD ["--version"]
