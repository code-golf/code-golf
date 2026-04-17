FROM alpine:3.23 AS builder

RUN apk add --no-cache clang cmake make

ENV LDFLAGS='-static' VER=4.5.0

WORKDIR /arkscript

ADD https://github.com/ArkScript-lang/Ark.git#v$VER .

RUN cmake -B build                  \
    -DARK_BUILD_EXE='ON'            \
    -DARK_ENABLE_SYSTEM='OFF'       \
    -DARK_STATIC='ON'               \
    -DCMAKE_BUILD_TYPE='MinSizeRel' \
    -DCMAKE_C_COMPILER='clang'      \
    -DCMAKE_CXX_COMPILER='clang++'  \
    -DCMAKE_INSTALL_PREFIX='/usr'   \
 && cmake --build build --parallel  \
 && cmake --install build --strip

FROM codegolf/lang-base

COPY --from=0 --parents /usr/bin/arkscript /usr/lib/Ark /

ENTRYPOINT ["arkscript"]

CMD ["--version"]
