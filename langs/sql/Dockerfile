FROM alpine:3.16 as builder

RUN apk add --no-cache curl gcc musl-dev

RUN curl https://www.sqlite.org/2022/sqlite-autoconf-3390000.tar.gz \
  | tar xz --strip 1

COPY sql.c /

RUN gcc -DSQLITE_ENABLE_MATH_FUNCTIONS -DSQLITE_THREADSAFE=0 \
    -static -s -o sql sql.c sqlite3.c

FROM codegolf/lang-base

COPY --from=0 /sql /usr/bin/

ENTRYPOINT ["sql"]

CMD ["-v"]
