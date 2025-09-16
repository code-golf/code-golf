FROM alpine:3.22 AS builder

# RUN apk add --no-cache clang cmake curl git libedit-dev libxml2-dev linux-headers lld ninja python3 zlib-dev

RUN apk add --no-cache clang cmake curl git libedit-static libxml2-static lld ninja python3 zlib-static

ENV CRT='x86_64-alpine-linux-musl' VER=21.1.1

RUN curl -#L https://github.com/llvm/llvm-project/archive/refs/tags/llvmorg-$VER.tar.gz \
  | tar xz --strip-components 1

WORKDIR /llvm/build

RUN cmake                              \
    -DCMAKE_BUILD_TYPE='Release'       \
    -DCMAKE_C_COMPILER='clang'         \
    -DCMAKE_CXX_COMPILER='clang++'     \
    -DCMAKE_CXX_STANDARD='20'          \
    -DCMAKE_GENERATOR='Ninja'          \
    -DCMAKE_INSTALL_PREFIX='/usr'      \
    -DLLVM_BUILD_STATIC='ON'           \
    -DLLVM_ENABLE_BINDINGS='OFF'       \
    -DLLVM_ENABLE_FFI='ON'             \
    -DLLVM_ENABLE_LLD='ON'             \
    -DLLVM_ENABLE_UNWIND_TABLES='OFF'  \
    -DLLVM_INCLUDE_BENCHMARKS='OFF'    \
    -DLLVM_INCLUDE_EXAMPLES='OFF'      \
    # -DLLVM_INCLUDE_TESTS='OFF'         \
    -DLLVM_HOST_TRIPLE=$CRT            \
    -DLLVM_PARALLEL_LINK_JOBS='2'      \
    # -DLLVM_STATIC_LINK_CXX_STDLIB='ON' \
    -DLLVM_TARGET_TRIPLE=$CRT          \
    -DLLVM_TARGETS_TO_BUILD='Native' ..

RUN ninja install

# RUN ldd /usr/bin/llc && ldd /usr/bin/lli && false

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/lli             /usr/bin/llvm-ir
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/

ENTRYPOINT ["llvm-ir"]

CMD ["--version"]
