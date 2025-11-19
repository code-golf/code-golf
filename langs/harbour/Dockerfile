FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev patch

ENV HB_BUILD_CONTRIBS='no' HB_BUILD_STRIP='bin' HB_USER_LDFLAGS='-static' VER=3.0.0

RUN curl -#L https://downloads.sourceforge.net/harbour-project/harbour-$VER.tar.xz \
  | tar xJ --strip-components 1

COPY harbour.patch /

RUN patch -p0 < harbour.patch

RUN ln -s /src/rtl/arc4.h /include/hbarc4.h

RUN make -j`nproc` install

COPY harbour.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/harbour -s -static harbour.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/harbour     /usr/bin/
COPY --from=0 /usr/local/bin/hbmk2 /usr/local/bin/

ENTRYPOINT ["harbour"]

CMD ["--version"]
