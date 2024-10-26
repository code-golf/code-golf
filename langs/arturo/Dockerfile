FROM alpine:3.20 AS builder

RUN apk add --no-cache curl g++ nim pcre

RUN curl -L https://github.com/arturo-lang/arturo/archive/0756a55.tar.gz \
  | tar xz

RUN cd arturo-*  \
 && ./build.nims \
    --install    \
    --mode mini  \
    --release

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /root/.arturo/bin/arturo /usr/bin/
COPY --from=0 /usr/lib/libpcre.so.1    /usr/lib/

ENTRYPOINT ["arturo"]

CMD ["--version"]
