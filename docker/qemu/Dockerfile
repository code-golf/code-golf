FROM alpine:3.21 AS builder

RUN apk add --no-cache curl build-base bash python3 samurai zlib-dev zlib-static gettext-dev gettext-static glib-dev glib-static

WORKDIR /root/src
ARG VER
RUN curl https://download.qemu.org/qemu-$VER.tar.xz -o qemu.tar.xz
RUN tar xJf qemu.tar.xz
WORKDIR ./qemu-$VER

COPY dump-state.patch .
RUN patch -i dump-state.patch -p0

COPY pid1.patch .
RUN patch -i pid1.patch -p0

# TODO: change x86-version to 4 for use on codegolf machine
ARG TARGETS
RUN ./configure --prefix=/usr --static --without-default-features --enable-linux-user \
 --target-list=${TARGETS// /,} \
 --enable-lto --enable-membarrier --disable-pie --enable-strip --x86-version=2 --disable-containers --disable-relocatable \
 && make -j`nproc` install

WORKDIR /root/as-ld-qemu
COPY as-ld-qemu.c .
COPY --from=asm asm.c asm.h .

RUN gcc -Wall -Werror -Wextra -pedantic -o /usr/bin/as-ld-qemu -s -static as-ld-qemu.c asm.c

RUN mkdir -p /root/dist; mv /usr/bin/as-ld-qemu /root/dist; for target in ${TARGETS}; do mv /usr/bin/qemu-${target//-linux-user/} /root/dist; done

FROM scratch AS export
COPY --from=builder /root/dist /usr/bin
