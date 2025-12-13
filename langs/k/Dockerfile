FROM alpine:3.23 AS builder

RUN apk add --no-cache curl gcc make musl-dev patch

ENV CC='gcc -static' VER=717063f249

WORKDIR /usr/bin

RUN curl -#L https://codeberg.org/ngn/k/archive/$VER.tar.gz \
  | tar xz --strip-components 1

COPY k.patch .

RUN patch -p0 < k.patch | make

FROM codegolf/lang-base

COPY --from=0 /usr/bin/k /usr/bin/

ENTRYPOINT ["k"]
