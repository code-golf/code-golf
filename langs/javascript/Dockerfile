FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl g++ git pkg-config procps python3

ENV PATH=/depot_tools:$PATH VER=13.6.233.10

# FIXME Pin to 3b73c8d, df1785b or later seems to fail with:
# lockfile.LockError: Error locking /depot_tools/external_bin/gsutil/gsutil_4.68
# (err: [Errno 11] Resource temporarily unavailable)
RUN git clone https://chromium.googlesource.com/chromium/tools/depot_tools.git \
 && (cd depot_tools && git reset --hard 3b73c8d) \
 && fetch v8

WORKDIR /v8

RUN git checkout $VER \
 && gclient sync -D

RUN gn gen out --args='                 \
    is_debug=false                      \
    symbol_level=0                      \
    use_thin_lto=true                   \
    v8_enable_gdbjit=false              \
    v8_enable_webassembly=false         \
    v8_monolithic=true                  \
    v8_use_external_startup_data=false' \
 && ninja -C out d8                     \
 && strip out/d8

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /v8/out/icudtl.dat                    /usr/bin/
COPY --from=0 /v8/out/d8                            /usr/bin/javascript

ENTRYPOINT ["javascript"]
