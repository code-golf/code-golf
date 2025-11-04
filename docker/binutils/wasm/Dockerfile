FROM emscripten/emsdk:4.0.18 AS builder
#FROM ghcr.io/webassembly/wasi-sdk AS builder

RUN apt-get update && \
    apt-get install -y git bison flex texinfo parallel --no-install-recommends && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root/src
ARG VER

RUN curl -#L https://ftpmirror.gnu.org/gnu/binutils/binutils-$VER.tar.xz -o binutils.tar.xz

RUN tar xJf binutils.tar.xz
WORKDIR /root/src/binutils-$VER

RUN sed -i '/^development=/s/true/false/' bfd/development.sh

ENV CFLAGS="-DHAVE_PSIGNAL=1 -DELIDE_CODE -Os -flto"
# -sDECLARE_ASM_MODULE_EXPORTS=0 is there to prevent minification of import names
ENV LDFLAGS="-flto -sDECLARE_ASM_MODULE_EXPORTS=0 -sFILESYSTEM=1 -sDETERMINISTIC=1 -sDYNAMIC_EXECUTION=0 -sALLOW_MEMORY_GROWTH=1 -sINITIAL_HEAP=65536"
#ENV LDFLAGS="-flto"

ARG TARGETS
RUN for target in ${TARGETS}; do mkdir -p /root/build/$target; cd /root/build/$target; emconfigure /root/src/binutils-$VER/configure --target="$target" --host=wasm32 --enable-default-execstack=no --enable-deterministic-archives --enable-ld=default --enable-new-dtags --disable-doc --disable-gprof --disable-nls --enable-binutils --disable-gdb --disable-gdbserver --disable-gprofng --disable-libdecnumber --disable-readline --disable-sim --disable-werror --disable-error-handling-script; done

COPY debug-info-for-all-sections.patch .
RUN patch -p1 < debug-info-for-all-sections.patch

COPY multiple-debug-info-per-line.patch .
RUN patch -p1 < multiple-debug-info-per-line.patch

COPY debug-info-for-directives.patch .
RUN patch -p1 < debug-info-for-directives.patch

RUN sed -i 's/#ifndef DWARF2_LINE_MIN_INSN_LENGTH/#if 1\n#undef DWARF2_LINE_MIN_INSN_LENGTH/' gas/dwarf2dbg.c

RUN for target in ${TARGETS}; do cd /root/build/$target; emmake make -j`nproc`; done

RUN mkdir -p /root/wasm; for target in ${TARGETS}; do \
    mv /root/build/$target/gas/as-new.wasm /root/wasm/${target}-as.wasm; \
    mv /root/build/$target/ld/ld-new.wasm /root/wasm/${target}-ld.wasm; \
    mv /root/build/$target/binutils/objdump.wasm /root/wasm/${target}-objdump.wasm; \
    done

FROM scratch AS export
COPY --from=builder /root/wasm /wasm
