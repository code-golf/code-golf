FROM alpine:3.22 AS builder

RUN apk add --no-cache cmake curl g++ linux-headers ninja python3

ENV VER=21.1.1

RUN curl -#L https://github.com/llvm/llvm-project/archive/refs/tags/llvmorg-$VER.tar.gz \
  | tar xz --strip-components 1

WORKDIR /llvm/build

RUN cmake                             \
    -DCMAKE_BUILD_TYPE='MinSizeRel'   \
    -DCMAKE_GENERATOR='Ninja'         \
    -DCMAKE_INSTALL_PREFIX='/usr'     \
    -DLLVM_ENABLE_BINDINGS='OFF'      \
    -DLLVM_ENABLE_UNWIND_TABLES='OFF' \
    -DLLVM_INCLUDE_BENCHMARKS='OFF'   \
    -DLLVM_INCLUDE_EXAMPLES='OFF'     \
    -DLLVM_INCLUDE_TESTS='OFF'        \
    -DLLVM_TARGETS_TO_BUILD='Native' ..

RUN ninja install \
 && strip /usr/bin/lli

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/lli             /usr/bin/llvm-ir
COPY --from=0 /usr/lib/libgcc_s.so.1   \
              /usr/lib/libstdc++.so.6  /usr/lib/

ENTRYPOINT ["llvm-ir"]

CMD ["--version"]
