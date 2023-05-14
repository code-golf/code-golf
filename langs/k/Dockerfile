FROM alpine:3.18 AS builder

ENV VERSION=c8f24002b1

RUN apk add --no-cache build-base curl

RUN curl -L https://codeberg.org/ngn/k/archive/$VERSION.tar.gz | tar xz \
 && sed -i s/march=native/march=x86-64-v3/ k/makefile                   \
 && make -C k CC="gcc -static"

COPY kwrapper.c /

RUN gcc -static -s -o kwrapper kwrapper.c

FROM codegolf/lang-base

COPY --from=0 /k/k /kwrapper /usr/bin/

ENTRYPOINT ["kwrapper"]
