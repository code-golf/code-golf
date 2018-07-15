FROM alpine:edge as builder

# Copied from https://github.com/AlexMasterov/dockerfiles/blob/master/alpine-libv8/6.9/Dockerfile
# I can't thank @AlexMasterov++ enough! V8 is a pig to build :-(

ARG BUILD_COMMIT=59ac9bda0b8fd61f7bc016be2c0cf58f8c3f714d
ARG BUILDTOOLS_COMMIT=5941c1b3df96c1db756a2834343533335c394c4a
ARG ICU_COMMIT=f61e46dbee9d539a32551493e3bcc1dea92f83ec
ARG GTEST_COMMIT=ce468a17c434e4e79724396ee1b51d86bfc8a88b
ARG TRACE_EVENT_COMMIT=211b3ed9d0481b4caddbee1322321b86a483ca1f
ARG GYP_COMMIT=d61a9397e668fa9843c4aa7da9e79460fe590bfb
ARG CLANG_COMMIT=ec200e7a8308f301ade483408bcc79ca17ceb11e
ARG JINJA2_COMMIT=45571de473282bd1d8b63a8dfcb1fd268d0635d2
ARG MARKUPSAFE_COMMIT=8f45f5cfa0009d2a70589bcda0349b8cb2b72783

ARG GN_SOURCE=https://www.dropbox.com/s/3ublwqh4h9dit9t/alpine-gn-80e00be.tar.gz
ARG V8_SOURCE=https://chromium.googlesource.com/v8/v8/+archive/6.9.231.tar.gz

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
