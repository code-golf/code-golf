FROM alpine:3.14 as builder

RUN mkdir /empty

RUN apk add --no-cache curl g++ libexecinfo-dev \
    linux-headers m4 make patch tar xz

WORKDIR /src

RUN curl ftp://ftp.gnu.org/gnu/gcc/gcc-11.2.0/gcc-11.2.0.tar.xz      | tar xJ \
 && curl ftp://gcc.gnu.org/pub/gcc/infrastructure/gmp-6.1.0.tar.bz2  | tar xj \
 && curl ftp://gcc.gnu.org/pub/gcc/infrastructure/mpfr-3.1.6.tar.bz2 | tar xj \
 && curl ftp://gcc.gnu.org/pub/gcc/infrastructure/mpc-1.0.3.tar.gz   | tar xz

RUN cd /src/gmp-6.1.0 \
 && ./configure --disable-shared --enable-static --prefix=/tmp/gcc \
 && make -j`nproc` install

RUN cd /src/mpfr-3.1.6 \
 && ./configure --disable-shared --enable-static --prefix=/tmp/gcc --with-gmp=/tmp/gcc \
 && make -j`nproc` install

RUN cd /src/mpc-1.0.3 \
 && ./configure --disable-shared --enable-static --prefix=/tmp/gcc --with-gmp=/tmp/gcc --with-mpfr=/tmp/gcc \
 && make -j`nproc` install

WORKDIR /scratch

RUN /src/gcc-11.2.0/configure     \
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
    /usr/local/libexec/gcc/x86_64-pc-linux-musl/11.2.0/f951

COPY fortran.c /scratch/

RUN gcc -s -o fortran fortran.c

FROM scratch

COPY --from=0 /scratch/fortran              /usr/bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1      \
              /lib/libz.so.1                /lib/
COPY --from=0 /usr/bin/as /usr/bin/ld       /usr/bin/
COPY --from=0 /usr/lib/crt1.o               \
              /usr/lib/crti.o               \
              /usr/lib/crtn.o               \
              /usr/lib/libbfd-2.35.2.so     \
              /usr/lib/libc.so              \
              /usr/lib/libctf.so.0          \
              /usr/lib/libm.a               \
              /usr/lib/libopcodes-2.35.2.so /usr/lib/
COPY --from=0 /usr/local/bin/gfortran       /usr/local/bin/
COPY --from=0 /usr/local/lib                /usr/local/lib
COPY --from=0 /usr/local/lib64              /usr/local/lib64
COPY --from=0 /usr/local/libexec/gcc/x86_64-pc-linux-musl/11.2.0/f951 \
              /usr/local/libexec/gcc/x86_64-pc-linux-musl/11.2.0/
COPY --from=0 /empty                        /proc
COPY --from=0 /empty                        /tmp

ENTRYPOINT ["/usr/bin/fortran"]

CMD ["--version"]
