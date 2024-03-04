FROM alpine:3.19 AS builder

RUN apk add --no-cache build-base cmake curl python3

ENV VER 17.0.6

RUN curl -L https://github.com/llvm/llvm-project/releases/download/llvmorg-$VER/llvm-project-$VER.src.tar.xz | tar xJ

WORKDIR /llvm-project-$VER.src/llvm/build

RUN cmake                            \
    -DCMAKE_INSTALL_PREFIX="/llvm"   \
    -DCMAKE_BUILD_TYPE="MinSizeRel"  \
    -DLLVM_TARGETS_TO_BUILD="host"   \
    -DLLVM_BUILD_TOOLS="NO"          \
    -DLLVM_BUILD_LLVM_DYLIB="YES"    \
    -DLLVM_INCLUDE_EXAMPLES="NO"     \
    -DLLVM_INCLUDE_TESTS="NO"        \
    -DLLVM_INCLUDE_BENCHMARKS="NO"   \
    -DLLVM_ENABLE_BINDINGS="NO"      \
    -DLLVM_ENABLE_UNWIND_TABLES="NO" \
    .. \
    && make -j`nproc` install

WORKDIR /llvm-project-$VER.src/clang/build

RUN cmake                            \
    -DCMAKE_INSTALL_PREFIX="/llvm"   \
    -DCMAKE_BUILD_TYPE="MinSizeRel"  \
    -DLLVM_TARGETS_TO_BUILD="host"   \
    -DLLVM_INCLUDE_TESTS="NO"        \
    -DLLVM_ENABLE_UNWIND_TABLES="NO" \
    .. \
    && make -j`nproc` install

RUN strip /llvm/bin/clang /llvm/lib/libLLVM-17.so

WORKDIR /runner

COPY cpp.c /runner/
RUN gcc -s -o cpp cpp.c

FROM codegolf/lang-base

COPY --from=0 /llvm/bin/clang          \
              /usr/bin/ld              \
              /runner/cpp              /usr/bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/
COPY --from=0 /usr/include/            /usr/include/
COPY --from=0 /llvm/lib/libLLVM-17.so  \
              /usr/lib/crt1.o          \
              /usr/lib/crti.o          \
              /usr/lib/crtn.o          \
              /usr/lib/libatomic.so    \
              /usr/lib/libbfd-2.41.so  \
              /usr/lib/libc.so         \
              /usr/lib/libctf.so.0     \
              /usr/lib/libffi.so.8     \
              /usr/lib/libgcc_s.so     \
              /usr/lib/libgcc_s.so.1   \
              /usr/lib/libjansson.so.4 \
              /usr/lib/liblzma.so.5    \
              /usr/lib/libsframe.so.1  \
              /usr/lib/libstdc++.so    \
              /usr/lib/libstdc++.so.6  \
              /usr/lib/libzstd.so.1    \
              /usr/lib/Scrt1.o         /usr/lib/

COPY --from=0 /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/libgcc.a    \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtbegin.o  \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtbeginS.o \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtend.o    \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/crtendS.o   \
              /usr/lib/gcc/x86_64-alpine-linux-musl/13.2.1/

COPY unbuffered.cpp /

ENTRYPOINT ["cpp"]

CMD ["--version"]
