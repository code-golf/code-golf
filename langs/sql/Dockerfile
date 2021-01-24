FROM alpine:3.13 as builder

RUN mkdir /empty

RUN apk add --no-cache curl gcc musl-dev

RUN curl https://www.sqlite.org/2021/sqlite-autoconf-3340100.tar.gz \
  | tar --strip 1 -xzf -

COPY sql.c /

RUN gcc -static -s -o sql sql.c sqlite3.c

FROM scratch

COPY --from=0 /empty /proc
COPY --from=0 /empty /tmp
COPY --from=0 /sql   /usr/bin/

ENTRYPOINT ["sql", "-v"]
