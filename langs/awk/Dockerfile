FROM alpine:3.16 as builder

RUN apk add --no-cache build-base curl

RUN curl https://ftp.gnu.org/gnu/gawk/gawk-5.2.0.tar.xz | tar xJ

RUN mv gawk-* gawk      \
 && cd gawk             \
 && ./configure         \
   --disable-lint       \
   --prefix /usr        \
 && make install        \
 && strip /usr/bin/gawk

FROM codegolf/lang-base

COPY --from=0 /usr/bin/gawk            /usr/bin/gawk
COPY --from=0 /usr/lib/gawk/           /usr/lib/gawk/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/

ENTRYPOINT ["gawk"]

CMD ["--version"]
