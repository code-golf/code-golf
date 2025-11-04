FROM alpine:3.21 AS builder

RUN apk add --no-cache curl build-base linux-headers gettext-dev gettext-static

WORKDIR /root/src
ARG VER
ADD https://sourceware.org/pub/binutils/releases/binutils-$VER.tar.xz binutils.tar.xz
RUN tar xJf binutils.tar.xz
WORKDIR ./binutils-$VER

ENV CFLAGS="-fno-PIE -no-pie -static -static-libgcc -s -Wl,-static" LDFLAGS="-static"

ARG TARGETS
RUN for target in ${TARGETS}; do mkdir -p /root/build/$target; cd /root/build/$target; /root/src/binutils-$VER/configure --target=$target --prefix=/usr --enable-gold=no --enable-ld=yes --disable-gdb --disable-gdbserver --disable-gprofng --disable-nls --enable-static --disable-shared --disable-multilib \
--enable-year2038 --disable-libquadmath --disable-libstdcxx --enable-lto --disable-bootstrap --disable-host-pie --disable-host-shared --without-pic; done

RUN for target in ${TARGETS}; do cd /root/build/$target; make -j`nproc`; done

RUN for target in ${TARGETS}; do cd /root/build/$target; make install; done

RUN mkdir -p /root/dist; for target in ${TARGETS}; do for tool in as ld objcopy; do mv /usr/$target/bin/$tool /root/dist/${target}-$tool; done; done

FROM scratch AS export
COPY --from=builder /root/dist /usr/bin
