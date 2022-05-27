FROM alpine:3.16 as builder

ENV VERSION=7508b9ae11

RUN apk add --no-cache build-base curl

RUN curl -L https://codeberg.org/ngn/k/archive/$VERSION.tar.gz | tar xz \
 && make -C k CC="gcc -static"

COPY kwrapper.c /

RUN gcc -static -s -o kwrapper kwrapper.c

FROM codegolf/lang-base

COPY --from=0 /k/k /kwrapper /usr/bin/

ENTRYPOINT ["kwrapper"]
