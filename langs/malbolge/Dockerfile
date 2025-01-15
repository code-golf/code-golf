FROM alpine:3.21 AS builder

RUN apk add --no-cache curl gcc musl-dev

ENV VER=07dce2d

RUN curl -#L https://github.com/graue/esofiles/tarball/$VER \
  | tar xz --strip-components 3

COPY malbolgewrapper.c /

RUN gcc -o /usr/bin/malbolge -s -static malbolge.c \
 && gcc -o /usr/bin/malbolgewrapper -s -static malbolgewrapper.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/malbolge        \
              /usr/bin/malbolgewrapper /usr/bin/

ENTRYPOINT ["malbolgewrapper"]

CMD ["--version"]
