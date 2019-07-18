FROM alpine:3.10 as builder

# Copied from https://github.com/AlexMasterov/dockerfiles/blob/master/alpine-libv8/7.4/Dockerfile
# I can't thank @AlexMasterov++ enough! V8 is a pig to build :-(

ARG BUILD_COMMIT=c1ab94d375f10578b3d207eca8fe4fb35274efb7
ARG BUILDTOOLS_COMMIT=6fbda1b24c1893a893b17aa219b765b9e7c801d8
ARG ICU_COMMIT=07e7295d964399ee7bee16a3ac7ca5a053b2cf0a
ARG GTEST_COMMIT=879ac092fde0a19e1b3a61b2546b2a422b1528bc
ARG TRACE_EVENT_COMMIT=e31a1706337ccb9a658b37d29a018c81695c6518
ARG CLANG_COMMIT=3114fbc11f9644c54dd0a4cdbfa867bac50ff983
ARG JINJA2_COMMIT=b41863e42637544c2941b574c7877d3e1f663e25
ARG MARKUPSAFE_COMMIT=8f45f5cfa0009d2a70589bcda0349b8cb2b72783
ARG CATAPULT_COMMIT=b6cc5a6baf93cfa6feeb240eea75c454506b0c3c

ARG GN_SOURCE=https://www.dropbox.com/s/3ublwqh4h9dit9t/alpine-gn-80e00be.tar.gz
ARG V8_SOURCE=https://chromium.googlesource.com/v8/v8/+archive/7.4.51.tar.gz

RUN set -x \
  && apk add --update --virtual .v8-build-dependencies \
    at-spi2-core-dev \
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
