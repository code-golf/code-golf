FROM alpine:3.21 AS builder

ENV VERSION=544d014afd

RUN apk add --no-cache build-base curl

RUN curl -L https://codeberg.org/ngn/k/archive/$VERSION.tar.gz | tar xz \
 && sed -i s/march=native/march=x86-64-v3/ k/makefile                   \
 && make -C k CC="gcc -static"

FROM codegolf/lang-base

COPY --from=0 /k/k /usr/bin/

ENTRYPOINT ["k"]
