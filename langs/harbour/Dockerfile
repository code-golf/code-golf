FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV HB_BUILD_CONTRIBS='no' HB_BUILD_STRIP='bin' HB_USER_LDFLAGS='-static' VER=53f3d24

RUN curl -#L https://github.com/harbour/core/tarball/$VER \
  | tar xz --strip-components 1

RUN make -j`nproc` install

COPY harbour.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/harbour -s -static harbour.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/harbour           /usr/bin/
COPY --from=0 /usr/local/bin/hbmk2       /usr/local/bin/
COPY --from=0 /usr/local/include/harbour /usr/local/include/harbour

ENTRYPOINT ["harbour"]
