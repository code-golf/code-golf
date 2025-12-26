FROM alpine:3.23 AS builder

RUN apk add --no-cache cmake curl g++ ninja python3

ENV VER=21.1.8

RUN curl -#L https://github.com/llvm/llvm-project/archive/refs/tags/llvmorg-$VER.tar.gz \
  | tar xz --strip-components 1

WORKDIR /llvm/build

RUN cmake -G Ninja                    \
    -DCMAKE_BUILD_TYPE='MinSizeRel'   \
    -DCMAKE_C_COMPILER='gcc'          \
    -DCMAKE_CXX_COMPILER='g++'        \
    -DCMAKE_INSTALL_PREFIX='/usr'     \
    -DLLVM_ENABLE_BINDINGS='OFF'      \
    -DLLVM_ENABLE_UNWIND_TABLES='OFF' \
    -DLLVM_INCLUDE_BENCHMARKS='OFF'   \
    -DLLVM_INCLUDE_EXAMPLES='OFF'     \
    -DLLVM_INCLUDE_TESTS='OFF'        \
    -DLLVM_PARALLEL_LINK_JOBS='2'     \
    -DLLVM_TARGETS_TO_BUILD='Native' ..

RUN ninja install-lli \
 && strip /usr/bin/lli

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1     /lib/
COPY --from=0 /usr/bin/lli            /usr/bin/llvm-ir
COPY --from=0 /usr/lib/libgcc_s.so.1  \
              /usr/lib/libstdc++.so.6 /usr/lib/

ENTRYPOINT ["llvm-ir"]

CMD ["--version"]
