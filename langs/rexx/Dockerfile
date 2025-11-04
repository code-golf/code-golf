FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV VER=3.9.7

RUN curl -#L https://downloads.sourceforge.net/regina-rexx/regina-rexx-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure      \
    --disable-flists \
    --disable-gci    \
    --disable-screen \
    --enable-64bit   \
 && make install

COPY rexx.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/rexx -s rexx.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/rexx       /usr/bin/
COPY --from=0 /usr/local/bin/rexx /usr/local/bin/

ENTRYPOINT ["rexx"]

CMD ["--version"]
