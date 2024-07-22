FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base git nim pcre

RUN git clone https://github.com/arturo-lang/arturo.git

WORKDIR /arturo

RUN ./build.nims -m mini -i -l

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /arturo/bin/arturo        /usr/bin/
COPY --from=0 /usr/lib/libpcre.so.1     /usr/lib/

ENTRYPOINT ["arturo"]

CMD ["-v"]
