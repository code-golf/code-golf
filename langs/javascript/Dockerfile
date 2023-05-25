FROM debian:bullseye-slim AS builder

RUN apt-get update \
 && apt-get install -y curl g++ git pkg-config python python3

ENV PATH "/depot_tools:$PATH"

RUN git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git \
 && fetch v8

WORKDIR /v8

RUN git checkout 11.3.244.12 && gclient sync -D

RUN gn gen out --args="                 \
    is_debug=false                      \
    symbol_level=0                      \
    use_thin_lto=true                   \
    v8_enable_gdbjit=false              \
    v8_enable_webassembly=false         \
    v8_monolithic=true                  \
    v8_use_external_startup_data=false" \
 && ninja -C out d8                     \
 && strip out/d8

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /v8/out/d8 /v8/out/icudtl.dat         /usr/bin/

ENTRYPOINT ["d8"]

CMD ["-v"]
