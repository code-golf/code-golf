FROM alpine:3.21 AS builder

RUN apk add --no-cache cmake curl g++ linux-headers make

ENV VER=cdd6e81

RUN curl -#L https://github.com/niallang/Nial_Development/tarball/$VER \
  | tar xz --strip-components 1

RUN cp /usr/include/linux/fcntl.h /usr/include/sys/

WORKDIR /BuildCore/build

RUN cmake ../src && make \
 && mv nialcore /usr/bin

WORKDIR /BuildNial/pkgblder

RUN mkdir ../src \
 && nialcore -defs buildfromcore

WORKDIR /BuildNial/build

RUN cmake ../src && make   \
 && mv nial /usr/local/bin \
 && strip /usr/local/bin/nial

COPY nial.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/nial -s /nial.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/nial            /usr/bin/
COPY --from=0 /usr/local/bin/nial      /usr/local/bin/

ENTRYPOINT ["nial"]

CMD ["--version"]
