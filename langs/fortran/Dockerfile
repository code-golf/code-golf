FROM alpine:3.23 AS builder

RUN apk add --no-cache clang m4 make

ENV CC='clang' CXX='clang++' LDFLAGS='-static' \
    GCC_VER=15.2.0 GMP_VER=6.3.0 MPC_VER=1.4.1 MPFR_VER=4.2.2

ADD --unpack https://ftpmirror.gnu.org/gnu/gcc/gcc-$GCC_VER/gcc-$GCC_VER.tar.xz /
ADD --unpack https://ftpmirror.gnu.org/gnu/gmp/gmp-$GMP_VER.tar.bz2             /
ADD --unpack https://ftpmirror.gnu.org/gnu/mpc/mpc-$MPC_VER.tar.xz              /
ADD --unpack https://ftpmirror.gnu.org/gnu/mpfr/mpfr-$MPFR_VER.tar.bz2          /

WORKDIR /gmp-$GMP_VER

RUN ./configure       \
    --disable-shared  \
    --prefix=/tmp/gmp \
 && make -j`nproc` install

WORKDIR /mpfr-$MPFR_VER

RUN ./configure         \
    --disable-shared    \
    --prefix=/tmp/mpfr  \
    --with-gmp=/tmp/gmp \
 && make -j`nproc` install

WORKDIR /mpc-$MPC_VER

RUN ./configure           \
    --disable-shared      \
    --prefix=/tmp/mpc     \
    --with-gmp=/tmp/gmp   \
    --with-mpfr=/tmp/mpfr \
 && make -j`nproc` install

WORKDIR /gcc-$GCC_VER

RUN ./configure                      \
    --disable-bootstrap              \
    --disable-gcov                   \
    --disable-libgomp                \
    --disable-libstdcxx-pch          \
    --disable-lto                    \
    --disable-multilib               \
    --disable-nls                    \
    --disable-shared                 \
    --enable-languages=fortran       \
    --with-gmp=/tmp/gmp              \
    --with-mpc=/tmp/mpc              \
    --with-mpfr=/tmp/mpfr            \
    --with-static-standard-libraries \
 && make -j`nproc`                   \
 && make install                     \
 && strip /usr/local/bin/gfortran    \
          /usr/local/libexec/gcc/x86_64-pc-linux-musl/$GCC_VER/f951

COPY fortran.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/fortran -s -static /fortran.c

FROM codegolf/lang-base

COPY --from=0           /usr/lib/crt1.o                   \
                        /usr/lib/crti.o                   \
                        /usr/lib/crtn.o                   \
                        /usr/lib/libbfd-2.45.1.so         \
                        /usr/lib/libc.so                  \
                        /usr/lib/libctf.so.0              \
                        /usr/lib/libjansson.so.4          \
                        /usr/lib/libm.a                   \
                        /usr/lib/libsframe.so.2           \
                        /usr/lib/libz.so.1                \
                        /usr/lib/libzstd.so.1             /usr/lib/
COPY --from=0 --parents /lib/ld-musl-*.so.1               \
                        /usr/bin/as                       \
                        /usr/bin/fortran                  \
                        /usr/bin/ld                       \
                        /usr/local/bin/gfortran           \
                        /usr/local/lib/gcc/*/*/crtbegin.o \
                        /usr/local/lib/gcc/*/*/crtend.o   \
                        /usr/local/lib/gcc/*/*/libgcc.a   \
                        /usr/local/lib64/libgfortran.a    \
                        /usr/local/lib64/libgfortran.spec \
                        /usr/local/lib64/libquadmath.a    \
                        /usr/local/libexec/gcc/*/*/f951   /

ENTRYPOINT ["fortran"]

CMD ["--version"]
