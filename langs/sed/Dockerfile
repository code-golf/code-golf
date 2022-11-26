FROM alpine:3.17 as builder

RUN apk add --no-cache build-base curl

RUN curl https://ftp.gnu.org/gnu/sed/sed-4.9.tar.xz | tar xJ

RUN cd sed-4.9                                 \
 && ./configure --enable-lto LDFLAGS="-static" \
 && make                                       \
 && strip sed/sed

FROM codegolf/lang-base

COPY --from=0 /sed-4.9/sed/sed /usr/bin/

ENTRYPOINT ["sed"]

CMD ["--version"]
