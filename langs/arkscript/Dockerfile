FROM alpine:3.23 AS builder

RUN apk add --no-cache cmake g++ make patch

ENV LDFLAGS='-static' VER=4.3.0

WORKDIR /arkscript

ADD https://github.com/ArkScript-lang/Ark.git#v$VER .

COPY arkscript.patch .

RUN patch -p0 < arkscript.patch

RUN cmake -B build                  \
    -DARK_BUILD_EXE='ON'            \
    -DARK_ENABLE_SYSTEM='OFF'       \
    -DCMAKE_BUILD_TYPE='MinSizeRel' \
    -DCMAKE_C_COMPILER='gcc'        \
    -DCMAKE_CXX_COMPILER='g++'      \
    -DCMAKE_INSTALL_PREFIX='/usr'   \
 && cmake --build build --parallel  \
 && cmake --install build --strip

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/arkscript  /usr/bin/
COPY --from=0 /usr/lib/Ark        /usr/lib/Ark

ENTRYPOINT ["arkscript"]

CMD ["--version"]
