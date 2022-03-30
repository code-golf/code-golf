FROM alpine:3.15 as builder

RUN mkdir /empty

RUN apk add --no-cache curl gcc musl-dev

RUN curl https://www.sqlite.org/2022/sqlite-autoconf-3380200.tar.gz \
  | tar xz --strip 1

COPY sql.c /

RUN gcc -DSQLITE_ENABLE_MATH_FUNCTIONS -DSQLITE_THREADSAFE=0 \
    -static -s -o sql sql.c sqlite3.c

FROM scratch

COPY --from=0 /empty /dev
COPY --from=0 /empty /proc
COPY --from=0 /empty /tmp
COPY --from=0 /sql   /usr/bin/

ENTRYPOINT ["sql"]

CMD ["-v"]
