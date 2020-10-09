# Based on https://gist.github.com/tylerchr/15a74b05944cfb90729db6a51265b6c9

FROM debian:buster-slim as pre-builder

RUN apt-get update && apt-get install -y curl git python

ENV PATH "/tmp/depot_tools:$PATH"

WORKDIR /tmp

RUN git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git \
 && fetch v8                \
 && cd v8                   \
 && git checkout 8.6.395.10 \
 && gclient sync

FROM alpine:3.12 as builder

COPY --from=0 /tmp/v8 /tmp/v8

RUN apk add --no-cache alpine-sdk binutils-gold clang curl g++ gcc git \
    glib-dev icu-dev libstdc++ linux-headers llvm9 make ninja python2 tar xz

ENV PATH "/usr/lib/llvm9/bin:$PATH"

RUN cp -f /usr/bin/ld.gold /usr/bin/ld               \
 && git clone https://gn.googlesource.com/gn /tmp/gn \
 && cd /tmp/gn                                       \
 && python build/gen.py                              \
 && ninja -C out                                     \
 && cp -f /tmp/gn/out/gn /tmp/v8/buildtools/linux64/gn

RUN cd /tmp/v8                               \
 && tools/dev/v8gen.py x64.release --        \
    'is_cfi=false'                           \
    'is_clang=false'                         \
    'is_component_build=true'                \
    'is_official_build=true'                 \
    'strip_debug_info=true'                  \
    'symbol_level=0'                         \
    'target_cpu="x64"'                       \
    'target_os="linux"'                      \
    'treat_warnings_as_errors=false'         \
    'use_custom_libcxx=false'                \
    'use_gold=false'                         \
    'use_sysroot=false'                      \
    'v8_enable_future=true'                  \
    'v8_enable_gdbjit=false'                 \
    'v8_enable_i18n_support=false'           \
    'v8_static_library=true'                 \
    'v8_target_cpu="x64"'                    \
    'v8_use_external_startup_data=false'     \
 && ninja -C out.gn/x64.release -j`nproc` d8 \
 && strip /tmp/v8/out.gn/x64.release/d8 /tmp/v8/out.gn/x64.release/*.so

FROM scratch

COPY --from=1 /lib/ld-musl-x86_64.so.1                        /lib/
COPY --from=1 /usr/lib/libgcc_s.so.1                          \
              /usr/lib/libstdc++.so.6                         /usr/lib/
COPY --from=1 /tmp/v8/include                                 /v8/include
COPY --from=1 /tmp/v8/out.gn/x64.release/d8                   \
              /tmp/v8/out.gn/x64.release/libchrome_zlib.so    \
              /tmp/v8/out.gn/x64.release/libv8.so             \
              /tmp/v8/out.gn/x64.release/libv8_libbase.so     \
              /tmp/v8/out.gn/x64.release/libv8_libplatform.so /v8/lib/

ENTRYPOINT ["/v8/lib/d8", "-v"]
