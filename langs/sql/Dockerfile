FROM alpine:3.12 as builder

RUN mkdir /empty

RUN apk add --no-cache curl gcc musl-dev

RUN curl https://www.sqlite.org/2020/sqlite-autoconf-3330000.tar.gz \
  | tar --strip 1 -xzf -

COPY sql.c /

RUN gcc -static -s -o sql sql.c sqlite3.c

FROM scratch

COPY --from=0 /empty /proc
COPY --from=0 /empty /tmp
COPY --from=0 /sql   /usr/bin/

ENTRYPOINT ["/usr/bin/sql", "-v"]
