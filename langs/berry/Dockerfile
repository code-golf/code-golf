FROM alpine:3.19 AS builder

RUN apk add --no-cache curl gcc make musl-dev python3 readline-dev

RUN curl -L https://github.com/berry-lang/berry/archive/97e0fca.tar.gz \
  | tar xz && mv berry-* berry

COPY berry.c berry/default/

RUN make -C berry

FROM codegolf/lang-base-no-proc

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /berry/berry             /usr/bin/

ENTRYPOINT ["berry"]

CMD ["-v"]
