FROM alpine:3.21 AS builder

ENV VER=3.9.6

RUN apk add --no-cache build-base curl

RUN curl -L https://sourceforge.net/projects/regina-rexx/files/regina-rexx/$VER/regina-rexx-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure \
 && make        \
 && make install

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh         /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/local/bin/regina    /usr/bin/
COPY --from=0 /libregina.so.3          /usr/lib/

COPY /rexx /usr/bin/

ENTRYPOINT ["rexx"]

CMD ["--version"]
