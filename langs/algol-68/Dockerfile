FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV LDFLAGS='-static' VER=3.10.6

RUN curl -#L https://jmvdveer.home.xs4all.nl/algol68g-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure   \
    --prefix=/usr \
 && make install  \
 && strip /usr/bin/a68g

FROM codegolf/lang-base

COPY --from=0 /usr/bin/a68g /usr/bin/algol-68

ENTRYPOINT ["algol-68"]

CMD ["--version"]
