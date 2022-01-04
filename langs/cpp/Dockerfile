FROM alpine:3.15 as builder

RUN mkdir /empty

RUN apk add --no-cache bash build-base cmake curl file git make python3 tar xz zip

RUN curl -L https://github.com/llvm/llvm-project/releases/download/llvmorg-13.0.0/llvm-13.0.0.src.tar.xz  | tar xJ \
 && curl -L https://github.com/llvm/llvm-project/releases/download/llvmorg-13.0.0/clang-13.0.0.src.tar.xz | tar xJ

WORKDIR /llvm-13.0.0.src/build

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

WORKDIR /clang-13.0.0.src/build

RUN cmake                            \
    -DCMAKE_INSTALL_PREFIX="/llvm"   \
    -DCMAKE_BUILD_TYPE="MinSizeRel"  \
    -DLLVM_TARGETS_TO_BUILD="host"   \
    -DLLVM_INCLUDE_TESTS="NO"        \
    -DLLVM_ENABLE_UNWIND_TABLES="NO" \
    .. \
    && make -j`nproc` install

RUN strip /llvm/bin/clang /llvm/lib/libLLVM-13.so

WORKDIR /runner

COPY cpp.c                           \
     unbuffered.cpp                  /runner/
RUN gcc -s -o cpp cpp.c
RUN g++ -c unbuffered.cpp

FROM scratch

COPY --from=0 /llvm/bin/clang          \
              /usr/bin/ld              \
              /runner/cpp              /usr/bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/
COPY --from=0 /usr/include/            /usr/include/
COPY --from=0 /llvm/lib/libLLVM-13.so  \
              /usr/lib/crt1.o          \
              /usr/lib/crti.o          \
              /usr/lib/crtn.o          \
              /usr/lib/libatomic.so    \
              /usr/lib/libbfd-2.37.so  \
              /usr/lib/libc.so         \
              /usr/lib/libctf.so.0     \
              /usr/lib/libffi.so.8     \
              /usr/lib/libgcc_s.so     \
              /usr/lib/libgcc_s.so.1   \
              /usr/lib/liblzma.so.5    \
              /usr/lib/libstdc++.so    \
              /usr/lib/libstdc++.so.6  \
              /usr/lib/Scrt1.o         /usr/lib/

COPY --from=0 usr/lib/gcc/x86_64-alpine-linux-musl/10.3.1/libgcc.a    \
              usr/lib/gcc/x86_64-alpine-linux-musl/10.3.1/crtbegin.o  \
              usr/lib/gcc/x86_64-alpine-linux-musl/10.3.1/crtbeginS.o \
              usr/lib/gcc/x86_64-alpine-linux-musl/10.3.1/crtend.o    \
              usr/lib/gcc/x86_64-alpine-linux-musl/10.3.1/crtendS.o   /usr/lib/gcc/x86_64-alpine-linux-musl/10.3.1/

COPY --from=0 /runner/unbuffered.o     /

COPY --from=0 /empty                   /tmp
COPY --from=0 /empty                   /proc

ENTRYPOINT ["/usr/bin/cpp"]

CMD ["--version"]
