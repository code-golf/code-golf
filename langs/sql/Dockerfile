FROM alpine:3.23 AS builder

RUN apk add --no-cache gcc musl-dev

ENV VER=3530200

ADD --unpack https://www.sqlite.org/2026/sqlite-autoconf-$VER.tar.gz /

WORKDIR /sqlite-autoconf-$VER

COPY sql.c .

RUN gcc -DSQLITE_ENABLE_MATH_FUNCTIONS -DSQLITE_THREADSAFE=0 \
    -o /usr/bin/sql -s -static sql.c sqlite3.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/sql /usr/bin/

ENTRYPOINT ["sql"]

CMD ["--version"]
