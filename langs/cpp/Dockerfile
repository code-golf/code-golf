FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base cmake curl linux-headers ninja python3

ENV VER=19.1.5

RUN curl -#L https://github.com/llvm/llvm-project/releases/download/llvmorg-$VER/llvm-project-$VER.src.tar.xz \
  | tar xJ --strip-components 1

WORKDIR /llvm/build

RUN cmake -G Ninja                    \
    -DCMAKE_BUILD_TYPE='MinSizeRel'   \
    -DCMAKE_INSTALL_PREFIX='/usr'     \
    -DLLVM_BUILD_LLVM_DYLIB='ON'      \
    -DLLVM_BUILD_TOOLS='OFF'          \
    -DLLVM_ENABLE_BINDINGS='OFF'      \
    -DLLVM_ENABLE_PROJECTS='clang'    \
    -DLLVM_ENABLE_UNWIND_TABLES='OFF' \
    -DLLVM_INCLUDE_BENCHMARKS='OFF'   \
    -DLLVM_INCLUDE_EXAMPLES='OFF'     \
    -DLLVM_INCLUDE_TESTS='OFF'        \
    -DLLVM_TARGETS_TO_BUILD='Native' ..

RUN ninja install        \
 && strip /usr/bin/clang \
          /usr/lib/libLLVM-19.so

COPY cpp.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/cpp -s -static /cpp.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/libz.so.1 /lib/
COPY --from=0 /usr/bin/clang /usr/bin/cpp /usr/bin/ld /usr/bin/
COPY --from=0 /usr/include/*.h                        /usr/include/
COPY --from=0 /usr/include/bits/*.h                   /usr/include/bits/
COPY --from=0 /usr/include/c++/13.2.1                 /usr/include/c++/13.2.1
COPY --from=0 /usr/lib/Scrt1.o                        \
              /usr/lib/crti.o                         \
              /usr/lib/crtn.o                         \
              /usr/lib/libLLVM-19.so                  \
              /usr/lib/libbfd-2.42.so                 \
              /usr/lib/libc.so                        \
              /usr/lib/libctf.so.0                    \
              /usr/lib/libgcc_s.so                    \
              /usr/lib/libgcc_s.so.1                  \
              /usr/lib/libjansson.so.4                \
              /usr/lib/libsframe.so.1                 \
              /usr/lib/libstdc++.so                   \
              /usr/lib/libstdc++.so.6                 \
              /usr/lib/libzstd.so.1                   /usr/lib/
COPY --from=0 /usr/lib/gcc/x86_64*/13.2.1/crtbegin.o  \
              /usr/lib/gcc/x86_64*/13.2.1/crtbeginS.o \
              /usr/lib/gcc/x86_64*/13.2.1/crtendS.o   \
              /usr/lib/gcc/x86_64*/13.2.1/libgcc.a    /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/

COPY unbuffered.cpp /

ENTRYPOINT ["cpp"]

CMD ["--version"]
