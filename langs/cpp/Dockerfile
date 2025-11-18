FROM alpine:3.22 AS builder

RUN apk add --no-cache cmake curl g++ linux-headers ninja python3

ENV VER=21.1.5

RUN curl -#L https://github.com/llvm/llvm-project/archive/refs/tags/llvmorg-$VER.tar.gz \
  | tar xz --strip-components 1

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
    -DLLVM_LINK_LLVM_DYLIB='ON'       \
    -DLLVM_TARGETS_TO_BUILD='Native' ..

RUN ninja install

RUN strip /usr/lib/libLLVM.so.21.1 \
          /usr/lib/libclang-cpp.so.21.1

COPY cpp.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/cpp -s -static /cpp.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1               /lib/
COPY --from=0 /usr/bin/clang                    \
              /usr/bin/clang++                  \
              /usr/bin/cpp                      \
              /usr/bin/ld                       /usr/bin/
COPY --from=0 /usr/include/*.h                  /usr/include/
COPY --from=0 /usr/include/bits/*.h             /usr/include/bits/
COPY --from=0 /usr/include/c++/14.2.0           /usr/include/c++/14.2.0
COPY --from=0 /usr/include/sys/*.h              /usr/include/sys/
COPY --from=0 /usr/lib/Scrt1.o                  \
              /usr/lib/crti.o                   \
              /usr/lib/crtn.o                   \
              /usr/lib/libLLVM.so.21.1          \
              /usr/lib/libbfd-2.44.so           \
              /usr/lib/libc.so                  \
              /usr/lib/libclang-cpp.so.21.1     \
              /usr/lib/libctf.so.0              \
              /usr/lib/libgcc_s.so              \
              /usr/lib/libgcc_s.so.1            \
              /usr/lib/libjansson.so.4          \
              /usr/lib/libsframe.so.1           \
              /usr/lib/libstdc++.so             \
              /usr/lib/libstdc++.so.6           \
              /usr/lib/libz.so.1                \
              /usr/lib/libzstd.so.1             /usr/lib/
COPY --from=0 /usr/lib/gcc/*/14.2.0/crtbegin.o  \
              /usr/lib/gcc/*/14.2.0/crtbeginS.o \
              /usr/lib/gcc/*/14.2.0/crtendS.o   \
              /usr/lib/gcc/*/14.2.0/libgcc.a    /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/

COPY unbuffered.cpp /

ENTRYPOINT ["cpp"]

CMD ["--version"]
