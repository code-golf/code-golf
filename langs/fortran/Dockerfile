ARG VER=14.2.0

FROM alpine:3.21 AS builder

ARG VER

RUN apk add --no-cache curl g++ linux-headers m4 make patch tar xz

WORKDIR /src

RUN curl ftp://ftp.gnu.org/gnu/gcc/gcc-$VER/gcc-$VER.tar.xz          | tar xJ \
 && curl ftp://gcc.gnu.org/pub/gcc/infrastructure/gmp-6.2.1.tar.bz2  | tar xj \
 && curl ftp://gcc.gnu.org/pub/gcc/infrastructure/mpfr-4.1.0.tar.bz2 | tar xj \
 && curl ftp://gcc.gnu.org/pub/gcc/infrastructure/mpc-1.2.1.tar.gz   | tar xz

RUN cd /src/gmp-6.2.1 \
 && ./configure --disable-shared --enable-static --prefix=/tmp/gcc \
 && make -j`nproc` install

RUN cd /src/mpfr-4.1.0 \
 && ./configure --disable-shared --enable-static --prefix=/tmp/gcc --with-gmp=/tmp/gcc \
 && make -j`nproc` install

RUN cd /src/mpc-1.2.1 \
 && ./configure --disable-shared --enable-static --prefix=/tmp/gcc --with-gmp=/tmp/gcc --with-mpfr=/tmp/gcc \
 && make -j`nproc` install

COPY fortran.c .

RUN gcc -o fortran fortran.c

WORKDIR /scratch

RUN /src/gcc-$VER/configure       \
    --disable-bootstrap           \
    --disable-gcov                \
    --disable-libgomp             \
    --disable-libstdcxx-pch       \
    --disable-lto                 \
    --disable-multilib            \
    --disable-nls                 \
    --disable-shared              \
    --enable-languages=fortran    \
    --enable-static               \
    --prefix=/usr/local           \
    --with-gmp=/tmp/gcc           \
    --with-mpc=/tmp/gcc           \
    --with-mpfr=/tmp/gcc          \
 && make -j`nproc`                \
 && make install                  \
 && strip /usr/local/bin/gfortran \
    /usr/local/libexec/gcc/x86_64-pc-linux-musl/$VER/f951

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/ld-musl-x86_64.so.1             /lib/
COPY --from=0 /src/fortran /usr/bin/as /usr/bin/ld /usr/bin/
COPY --from=0 /usr/lib/crt1.o                      \
              /usr/lib/crti.o                      \
              /usr/lib/crtn.o                      \
              /usr/lib/libbfd-2.43.1.so            \
              /usr/lib/libc.so                     \
              /usr/lib/libctf.so.0                 \
              /usr/lib/libjansson.so.4             \
              /usr/lib/libm.a                      \
              /usr/lib/libopcodes-2.43.1.so        \
              /usr/lib/libsframe.so.1              \
              /usr/lib/libz.so.1                   \
              /usr/lib/libzstd.so.1                /usr/lib/
COPY --from=0 /usr/local/bin/gfortran              /usr/local/bin/
COPY --from=0 /usr/local/lib                       /usr/local/lib
COPY --from=0 /usr/local/lib64                     /usr/local/lib64
COPY --from=0 /usr/local/libexec/gcc/x86_64-pc-linux-musl/$VER/f951 \
              /usr/local/libexec/gcc/x86_64-pc-linux-musl/$VER/

ENTRYPOINT ["fortran"]

CMD ["--version"]
