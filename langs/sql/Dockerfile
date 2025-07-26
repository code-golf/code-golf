FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc musl-dev

ENV VER=3500300

RUN curl -#L https://www.sqlite.org/2025/sqlite-autoconf-$VER.tar.gz \
  | tar xz --strip-components 1

COPY sql.c /

RUN gcc -DSQLITE_ENABLE_MATH_FUNCTIONS -DSQLITE_THREADSAFE=0 \
    -o /usr/bin/sql -s -static sql.c sqlite3.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/sql /usr/bin/

ENTRYPOINT ["sql"]

CMD ["--version"]
