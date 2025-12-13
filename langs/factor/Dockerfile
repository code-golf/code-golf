FROM alpine:3.22 AS builder

RUN apk add --no-cache bash clang curl git make patch

ENV VER=0.101

RUN curl -#o factor.zip https://downloads.factorcode.org/releases/$VER/factor-src-$VER.zip \
 && unzip -n factor

COPY factor.patch /

RUN patch -p0 < factor.patch

WORKDIR /factor

RUN ./build.sh compile \
 && ./build.sh bootstrap

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1        /lib/
COPY --from=0 /factor/factor             /usr/bin/
COPY --from=0 /factor/factor.image.fresh /usr/bin/factor.image
COPY --from=0 /usr/lib/libgcc_s.so.1     \
              /usr/lib/libstdc++.so.6    /usr/lib/
COPY --from=0 /factor/basis              /usr/lib/basis
COPY --from=0 /factor/core               /usr/lib/core

ENTRYPOINT ["factor"]

CMD ["--version"]
