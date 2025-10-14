FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl g++ git ninja-build python3

ENV PATH=/depot_tools:$PATH VER=3.9.4

RUN git clone --progress https://chromium.googlesource.com/chromium/tools/depot_tools \
 && fetch dart

WORKDIR /sdk

RUN git checkout $VER \
 && gclient sync -D

RUN gn gen out --args='         \
    dart_include_wasm_opt=false \
    dart_platform_sdk=true      \
    dart_runtime_mode="release" \
    dart_sdk_output=""          \
    dart_version_git_info=false \
    is_debug=false              \
    is_product=true             \
    target_cpu="x64"'           \
 && ninja -C out create_sdk     \
 && strip out/bin/dart

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                     \
              /lib/x86_64-linux-gnu/libm.so.6                     /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                         /lib64/
COPY --from=0 /sdk/out/bin/dart                                   /usr/bin/
COPY --from=0 /sdk/out/bin/snapshots/kernel-service.dart.snapshot /usr/bin/snapshots/
COPY --from=0 /sdk/out/include                                    /usr/include
COPY --from=0 /sdk/out/lib                                        /usr/lib

ENTRYPOINT ["dart"]

CMD ["--version"]
