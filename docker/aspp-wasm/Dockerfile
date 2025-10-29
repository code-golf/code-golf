FROM emscripten/emsdk:4.0.18 AS builder

WORKDIR /root/src/aspp
COPY aspp.c .
COPY --from=asm asm.c asm.h .

# -sDECLARE_ASM_MODULE_EXPORTS=0 is there to prevent minification of import names
ENV LDFLAGS="-flto -sDECLARE_ASM_MODULE_EXPORTS=0 -sFILESYSTEM=1 -sDETERMINISTIC=1 -sDYNAMIC_EXECUTION=0 -sALLOW_MEMORY_GROWTH=1 -sINITIAL_HEAP=65536"

RUN emcc -Wall -pedantic -O2 -flto aspp.c asm.c -o aspp

FROM scratch AS export
COPY --from=builder /root/src/aspp/aspp.wasm /wasm/
