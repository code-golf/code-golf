FROM alpine:3.23 AS builder

RUN apk add --no-cache clang cmake make

ENV LDFLAGS='-static' VER=2.3

WORKDIR /lily

ADD https://gitlab.com/FascinatedBox/lily.git#v$VER .

RUN cmake -B build                  \
    -DCMAKE_BUILD_TYPE='MinSizeRel' \
    -DCMAKE_C_COMPILER='clang'      \
    -DCMAKE_CXX_COMPILER='clang++'  \
    -DCMAKE_INSTALL_PREFIX='/usr'   \
 && cmake --build build --parallel  \
 && cmake --install build --strip

FROM codegolf/lang-base

COPY --from=0 /usr/bin/lily /usr/bin/

ENTRYPOINT ["lily"]
