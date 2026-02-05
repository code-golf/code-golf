FROM alpine:3.23 AS builder

RUN apk add --no-cache build-base cmake git

ENV VER=2.3

RUN git clone --branch v$VER https://gitlab.com/FascinatedBox/lily.git

WORKDIR /lily

RUN cmake                           \
    -DCMAKE_BUILD_TYPE='MinSizeRel' \
    -DCMAKE_INSTALL_PREFIX='/usr' . \
 && make install                    \
 && strip /usr/bin/lily

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/lily            /usr/bin/

ENTRYPOINT ["lily"]
