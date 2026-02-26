FROM alpine:3.23 AS builder

RUN apk add --no-cache cmake g++ make

ENV VER=4.3.0

ADD https://github.com/ArkScript-lang/Ark.git#v$VER /arkscript

# HACK Tweak the CMake file to link libArkReactor statically.
RUN sed -i 's/ArkReactor SHARED/ArkReactor STATIC/' /arkscript/CMakeLists.txt

RUN cmake -B build -S arkscript     \
    -DARK_BUILD_EXE='ON'            \
    -DARK_ENABLE_SYSTEM='OFF'       \
    -DCMAKE_BUILD_TYPE='MinSizeRel' \
    -DCMAKE_C_COMPILER='gcc'        \
    -DCMAKE_CXX_COMPILER='g++'      \
    -DCMAKE_INSTALL_PREFIX='/usr'   \
 && cmake --build build --parallel  \
 && cmake --install build --strip

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1     /lib/
COPY --from=0 /usr/bin/arkscript      /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1  \
              /usr/lib/libstdc++.so.6 /usr/lib/
COPY --from=0 /usr/lib/Ark/std/*.ark  /usr/lib/std/

ENTRYPOINT ["arkscript"]

CMD ["--version"]
