FROM alpine:3.22 AS builder

RUN apk add --no-cache coreutils

ENV VER=0.6.5

RUN wget -qO cjam.jar https://sourceforge.net/projects/cjam/files/cjam-$VER/cjam-$VER.jar/download

FROM codegolf/lang-java

COPY --from=0 /cjam.jar                   /
COPY --from=0 /usr/lib/libacl.so.1        \
              /usr/lib/libattr.so.1       \
              /usr/lib/libcrypto.so.3     \
              /usr/lib/libskarnet.so.2.14 \
              /usr/lib/libutmps.so.0.1    /usr/lib/

ENTRYPOINT ["cjam"]

CMD ["--version"]
