FROM alpine:3.23 AS builder

RUN apk add --no-cache clang cmake make patch

ENV LDFLAGS='-static' VER=4.3.2

WORKDIR /arkscript

ADD https://github.com/ArkScript-lang/Ark.git#v$VER .

COPY arkscript.patch .

RUN patch -p0 < arkscript.patch

RUN cmake -B build                  \
    -DARK_BUILD_EXE='ON'            \
    -DARK_ENABLE_SYSTEM='OFF'       \
    -DCMAKE_BUILD_TYPE='MinSizeRel' \
    -DCMAKE_C_COMPILER='clang'      \
    -DCMAKE_CXX_COMPILER='clang++'  \
    -DCMAKE_INSTALL_PREFIX='/usr'   \
 && cmake --build build --parallel  \
 && cmake --install build --strip

FROM codegolf/lang-base

COPY --from=0 /usr/bin/arkscript /usr/bin/
COPY --from=0 /usr/lib/Ark       /usr/lib/Ark

ENTRYPOINT ["arkscript"]

CMD ["--version"]
