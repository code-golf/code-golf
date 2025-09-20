FROM alpine:3.22 AS builder

RUN apk add --no-cache cmake curl gcc make musl-dev

ENV CFLAGS='-Wno-implicit-function-declaration' LDFLAGS='-static' VER=b7bf640

RUN curl -#L https://github.com/VorpalBlade/cfunge/tarball/$VER \
  | tar xz --strip-components 1

WORKDIR /build

RUN cmake                            \
    -DCMAKE_BUILD_TYPE='MinSizeRel'  \
    -DCMAKE_INSTALL_PREFIX='/usr' .. \
 && make install                     \
 && strip /usr/bin/cfunge

COPY befunge.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/befunge -s -static /befunge.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/befunge /usr/bin/cfunge /usr/bin/

ENTRYPOINT ["befunge"]

CMD ["--version"]
