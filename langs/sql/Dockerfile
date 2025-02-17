FROM alpine:3.21 AS builder

RUN apk add --no-cache curl gcc musl-dev

RUN curl https://www.sqlite.org/2025/sqlite-autoconf-3490000.tar.gz \
  | tar xz --strip 1

COPY sql.c /

RUN gcc -DSQLITE_ENABLE_MATH_FUNCTIONS -DSQLITE_THREADSAFE=0 \
    -static -s -o sql sql.c sqlite3.c

FROM codegolf/lang-base

COPY --from=0 /sql /usr/bin/

ENTRYPOINT ["sql"]

CMD ["--version"]
