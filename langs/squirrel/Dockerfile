FROM alpine:3.24 AS builder

RUN apk add --no-cache clang cmake gcc make patch

ENV VER=3.2

ADD https://github.com/albertodemichelis/squirrel.git#v$VER .

COPY squirrel.patch /

RUN patch -p0 < squirrel.patch

WORKDIR /squirrel

RUN cmake ..                             \
    -DCMAKE_BUILD_TYPE='MinSizeRel'      \
    -DCMAKE_POLICY_VERSION_MINIMUM='3.5' \
 && make install

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1      /lib/
COPY --from=0 /usr/local/bin/sq_static /usr/bin/squirrel
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/

ENTRYPOINT ["squirrel"]

CMD ["-v"]
