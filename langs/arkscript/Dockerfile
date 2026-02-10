FROM alpine:3.23 AS builder

RUN apk add --no-cache cmake g++ git make

ENV VER=4.2.0

# Use git rather than curl because GitHub archives don't contain submodules
# See https://github.com/dear-github/dear-github/issues/214
RUN git clone --branch v$VER --recursive \
    https://github.com/ArkScript-lang/Ark

RUN cmake -S Ark -B Ark             \
    -DARK_BUILD_EXE='ON'            \
    -DARK_ENABLE_SYSTEM='OFF'       \
    -DCMAKE_BUILD_TYPE='MinSizeRel' \
    -DCMAKE_C_COMPILER='gcc'        \
    -DCMAKE_CXX_COMPILER='g++'      \
 && cmake --build Ark --parallel    \
 && cmake --install Ark --strip

COPY arkscript.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/arkscript -s arkscript.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1             /lib/
COPY --from=0 /usr/bin/arkscript              /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1          \
              /usr/lib/libstdc++.so.6         /usr/lib/
COPY --from=0 /usr/local/bin/arkscript        /usr/local/bin/
COPY --from=0 /usr/local/lib/libArkReactor.so /usr/local/lib/
COPY --from=0 /usr/local/lib/Ark/std          /usr/local/lib/std

ENV PATH=/usr/bin:$PATH

ENTRYPOINT ["arkscript"]

CMD ["--version"]
