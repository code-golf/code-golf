FROM alpine:3.22 AS builder

RUN apk add --no-cache bash clang curl git make

ENV VER=0.100

RUN curl -#o factor.zip https://downloads.factorcode.org/releases/$VER/factor-src-$VER.zip \
 && unzip -n factor

WORKDIR /factor

RUN sed -i 's/64//' basis/io/directories/unix/linux/linux.factor   \
 && sed -i 's/lxstat64/lxstat/' basis/unix/stat/linux/64/64.factor \
 && sed -i '/^\[/,/al/d' basis/vocabs/prettyprint/prettyprint.factor

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
