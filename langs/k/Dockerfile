FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev patch

ENV CC='gcc -static' VER=b9eeb91e03

WORKDIR /k

ADD https://codeberg.org/ngn/k.git?commit=$VER .

COPY k.patch .

RUN patch -p0 < k.patch | make

FROM codegolf/lang-base

COPY --from=0 /k/k /usr/bin/

ENTRYPOINT ["k"]
