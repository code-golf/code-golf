FROM alpine:3.15 as builder

ENV VERSION=a52e48b43b

RUN mkdir /empty

RUN apk add --no-cache build-base curl

RUN curl -L https://codeberg.org/ngn/k/archive/$VERSION.tar.gz | tar xz \
 && cd /k                                                               \
 && make

COPY kwrapper.c /

RUN gcc -s -o kwrapper kwrapper.c

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1          /lib/
COPY --from=0 /empty                            /proc
COPY --from=0 /empty                            /tmp
COPY --from=0 /k/k                              \
              /kwrapper                         /usr/bin/

ENTRYPOINT ["kwrapper"]
