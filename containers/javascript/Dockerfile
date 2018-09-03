FROM alpine:edge as builder

# Copied from https://github.com/AlexMasterov/dockerfiles/blob/master/alpine-libv8/7.0/Dockerfile
# I can't thank @AlexMasterov++ enough! V8 is a pig to build :-(

ARG BUILD_COMMIT=03bbc870bf80775a21c613a19b041e78eb334c65
ARG BUILDTOOLS_COMMIT=2dff9c9c74e9d732e6fe57c84ef7fd044cc45d96
ARG ICU_COMMIT=a191af9d025859e8368b8b469120d78006e9f5f6
ARG GTEST_COMMIT=d5266326752f0a1dadbd310932d8f4fd8c3c5e7d
ARG TRACE_EVENT_COMMIT=211b3ed9d0481b4caddbee1322321b86a483ca1f
ARG GYP_COMMIT=d61a9397e668fa9843c4aa7da9e79460fe590bfb
ARG CLANG_COMMIT=bb4146fb8a9dde405b71914657bb461dc93912ab
ARG JINJA2_COMMIT=b41863e42637544c2941b574c7877d3e1f663e25
ARG MARKUPSAFE_COMMIT=8f45f5cfa0009d2a70589bcda0349b8cb2b72783

ARG GN_SOURCE=https://www.dropbox.com/s/3ublwqh4h9dit9t/alpine-gn-80e00be.tar.gz
ARG V8_SOURCE=https://chromium.googlesource.com/v8/v8/+archive/7.0.272.tar.gz

RUN set -x \
  && apk add --update --virtual .v8-build-dependencies \
    curl \
    g++ \
    gcc \
    glib-dev \
    icu-dev \
    linux-headers \
    make \
    ninja \
    python \
    tar \
    xz \
  && : "---------- V8 ----------" \
  && mkdir /v8 \
  && curl -fSL --connect-timeout 30 ${V8_SOURCE} | tar xmz -C /v8 \
  && : "---------- Dependencies ----------" \
  && DEPS=" \
    chromium/buildtools.git@${BUILDTOOLS_COMMIT}:buildtools; \
    chromium/src/build.git@${BUILD_COMMIT}:build; \
    chromium/src/base/trace_event/common.git@${TRACE_EVENT_COMMIT}:base/trace_event/common; \
    chromium/src/tools/clang.git@${CLANG_COMMIT}:tools/clang; \
    chromium/src/third_party/jinja2.git@${JINJA2_COMMIT}:third_party/jinja2; \
    chromium/src/third_party/markupsafe.git@${MARKUPSAFE_COMMIT}:third_party/markupsafe; \
    chromium/deps/icu.git@${ICU_COMMIT}:third_party/icu; \
    external/gyp.git@${GYP_COMMIT}:tools/gyp; \
    external/github.com/google/googletest.git@${GTEST_COMMIT}:third_party/googletest/src \
  " \
  && while [ "${DEPS}" ]; do \
    dep="${DEPS%%;*}" \
    link="${dep%%:*}" \
    url="${link%%@*}" url="${url#"${url%%[![:space:]]*}"}" \
    hash="${link#*@}" \
    dir="${dep#*:}"; \
    [ -n "${dep}" ] \
      && dep_url="https://chromium.googlesource.com/${url}/+archive/${hash}.tar.gz" \
      && dep_dir="/v8/${dir}" \
      && mkdir -p ${dep_dir} \
      && curl -fSL --connect-timeout 30 ${dep_url} | tar xmz -C ${dep_dir} \
      & [ "${DEPS}" = "${dep}" ] && DEPS='' || DEPS="${DEPS#*;}"; \
    done; \
    wait \
  && : "---------- Downloads the current stable Linux sysroot ----------" \
  && /v8/build/linux/sysroot_scripts/install-sysroot.py --arch=amd64 \
  && : "---------- Proper GN ----------" \
  && apk add --virtual .gn-runtime-dependencies \
    libevent \
    libexecinfo \
    libstdc++ \
  && curl -fSL --connect-timeout 30 ${GN_SOURCE} | tar xmz -C /v8/buildtools/linux64/ \
  && : "---------- Build instructions ----------" \
  && cd /v8 \
  && ./tools/dev/v8gen.py \
    x64.release \
    -- \
      binutils_path=\"/usr/bin\" \
      target_os=\"linux\" \
      target_cpu=\"x64\" \
      v8_target_cpu=\"x64\" \
      v8_use_external_startup_data=false \
      v8_enable_future=true \
      is_official_build=true \
      is_component_build=true \
      is_cfi=false \
      is_clang=false \
      use_custom_libcxx=false \
      use_sysroot=false \
      use_gold=false \
      use_allocator_shim=false \
      treat_warnings_as_errors=false \
      symbol_level=0 \
  && : "---------- Build ----------" \
  && ninja d8 -C out.gn/x64.release/ -j $(getconf _NPROCESSORS_ONLN)

RUN echo -e "#!/bin/sh -e\n\
\n\
/bin/cat - > /tmp/code.js\n\
\n\
exec /v8/out.gn/x64.release/d8 /tmp/code.js \"\$@\"" > /usr/bin/javascript \
 && chmod +x /usr/bin/javascript

FROM scratch

COPY --from=0 /bin/cat                                      \
              /bin/sh                                       /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1                      /lib/
COPY --from=0 /usr/bin/javascript                           /usr/bin/
COPY --from=0 /usr/lib/libgcc_s.so.1                        \
              /usr/lib/libstdc++.so.6                       /usr/lib/
COPY --from=0 /v8/out.gn/x64.release/d8                     /v8/out.gn/x64.release/
COPY --from=0 /v8/out.gn/x64.release/./libicui18n.so        \
              /v8/out.gn/x64.release/./libicuuc.so          \
              /v8/out.gn/x64.release/./libv8.so             \
              /v8/out.gn/x64.release/./libv8_libbase.so     \
              /v8/out.gn/x64.release/./libv8_libplatform.so /v8/out.gn/x64.release/./

ENTRYPOINT ["/v8/out.gn/x64.release/d8", "-v"]
