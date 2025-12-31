FROM alpine:3.23 AS builder

RUN apk add --no-cache clang cmake curl make

ENV CC='clang' CXX='clang++' CFLAGS='-Wno-cpp' CMAKE_POLICY_VERSION_MINIMUM='3.5' LDFLAGS='-static' VER=cdd6e81

RUN curl -#L https://github.com/niallang/Nial_Development/tarball/$VER \
  | tar xz --strip-components 2

WORKDIR /build

RUN cmake ../src && make \
 && mv nialcore /usr/bin

WORKDIR /pkgblder

RUN nialcore -defs buildfromcore

WORKDIR /build

RUN cmake ../src && make   \
 && mv nial /usr/local/bin \
 && strip /usr/local/bin/nial

COPY nial.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/nial -s -static /nial.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/nial       /usr/bin/
COPY --from=0 /usr/local/bin/nial /usr/local/bin/

ENTRYPOINT ["nial"]
