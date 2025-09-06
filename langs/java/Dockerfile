FROM alpine:3.22 AS builder

RUN apk add --no-cache alsa-lib-dev autoconf bash build-base cups-dev curl fontconfig-dev \
                       libxrandr-dev libxrender-dev libxt-dev libxtst-dev musl-locales zip

# Also rebuild CJam, Groovy, Kotlin, Scala, Vyxal.
ENV VER=24.0.2 BUILD=12

RUN curl -#L https://github.com/adoptium/temurin24-binaries/releases/download/jdk-$VER+$BUILD/OpenJDK24U-jdk_x64_alpine-linux_hotspot_${VER}_$BUILD.tar.gz \
  | tar xz --directory /tmp --strip-components 1

RUN curl -#L https://github.com/adoptium/temurin24-binaries/releases/download/jdk-$VER+$BUILD/OpenJDK24U-jdk-sources_${VER}_$BUILD.tar.gz \
  | tar xz --strip-components 1

RUN bash configure                     \
    --enable-headless-only             \
    --enable-jvm-feature-link-time-opt \
    --enable-jvm-feature-management    \
    --enable-jvm-feature-opt-size      \
    --with-boot-jdk=/tmp               \
    --with-conf-name=release           \
    --with-jobs=`nproc`                \
 && make images

RUN /build/release/images/jdk/bin/jlink           \
    --add-modules java.base,java.se               \
    --add-modules jdk.compiler,jdk.unsupported    \
    --compress zip-9                              \
    --module-path /build/release/images/jdk/jmods \
    --output /opt/java

COPY java.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/java -s -static java.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1              /lib/
COPY --from=0 /opt/java/bin/java                    /opt/java/bin/
COPY --from=0 /opt/java/conf/security/java.security /opt/java/conf/security/
COPY --from=0 /opt/java/lib/jrt-fs.jar              \
              /opt/java/lib/jvm.cfg                 \
              /opt/java/lib/libjava.so              \
              /opt/java/lib/libjimage.so            \
              /opt/java/lib/libjli.so               \
              /opt/java/lib/libmanagement.so        \
              /opt/java/lib/libnet.so               \
              /opt/java/lib/libnio.so               \
              /opt/java/lib/libverify.so            \
              /opt/java/lib/libzip.so               \
              /opt/java/lib/modules                 \
              /opt/java/lib/tzdb.dat                /opt/java/lib/
COPY --from=0 /opt/java/lib/server/libjvm.so        /opt/java/lib/server/
COPY --from=0 /usr/bin/java                         /usr/bin/
COPY --from=0 /usr/lib/libz.so.1                    /usr/lib/

ENTRYPOINT ["java"]

CMD ["--version"]
