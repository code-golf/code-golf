FROM alpine:3.15 as builder

ENV VERSION=cf753e859d

RUN mkdir /empty

RUN apk add --no-cache build-base curl

RUN curl -L https://codeberg.org/ngn/k/archive/$VERSION.tar.gz | tar xz \
 && make -C k CC="gcc -static"

COPY kwrapper.c /

RUN gcc -static -s -o kwrapper kwrapper.c

FROM scratch

COPY --from=0 /empty         /dev
COPY --from=0 /empty         /proc
COPY --from=0 /empty         /tmp
COPY --from=0 /k/k /kwrapper /usr/bin/

ENTRYPOINT ["kwrapper"]
