ARG VER=15.2.0

FROM alpine:3.22 AS builder

ARG VER

RUN apk add --no-cache curl g++ m4 make

ENV GMPVER=6.3.0 LDFLAGS='-static' MPCVER=1.3.1 MPFRVER=4.2.2

WORKDIR /gmp

RUN curl -#L https://ftpmirror.gnu.org/gnu/gmp/gmp-$GMPVER.tar.bz2 \
  | tar xj --strip-components 1

RUN ./configure       \
    --disable-shared  \
    --prefix=/tmp/gmp \
 && make -j`nproc` install

WORKDIR /mpfr

RUN curl -#L https://ftpmirror.gnu.org/gnu/mpfr/mpfr-$MPFRVER.tar.bz2 \
  | tar xj --strip-components 1

RUN ./configure         \
    --disable-shared    \
    --prefix=/tmp/mpfr  \
    --with-gmp=/tmp/gmp \
 && make -j`nproc` install

WORKDIR /mpc

RUN curl -#L https://ftpmirror.gnu.org/gnu/mpc/mpc-$MPCVER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure           \
    --disable-shared      \
    --prefix=/tmp/mpc     \
    --with-gmp=/tmp/gmp   \
    --with-mpfr=/tmp/mpfr \
 && make -j`nproc` install

WORKDIR /fortran

RUN curl -#L https://ftpmirror.gnu.org/gnu/gcc/gcc-$VER/gcc-$VER.tar.xz \
  | tar xJ --strip-components 1

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
          /usr/local/libexec/gcc/x86_64-pc-linux-musl/$VER/f951

COPY fortran.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/fortran -s -static /fortran.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/ld-musl-*.so.1               /lib/
COPY --from=0 /usr/bin/as                       \
              /usr/bin/fortran                  \
              /usr/bin/ld                       /usr/bin/
COPY --from=0 /usr/lib/crt1.o                   \
              /usr/lib/crti.o                   \
              /usr/lib/crtn.o                   \
              /usr/lib/libbfd-2.44.so           \
              /usr/lib/libc.so                  \
              /usr/lib/libctf.so.0              \
              /usr/lib/libjansson.so.4          \
              /usr/lib/libm.a                   \
              /usr/lib/libsframe.so.1           \
              /usr/lib/libz.so.1                \
              /usr/lib/libzstd.so.1             /usr/lib/
COPY --from=0 /usr/local/bin/gfortran           /usr/local/bin/
COPY --from=0 /usr/local/lib/gcc/*/*/crtbegin.o \
              /usr/local/lib/gcc/*/*/crtend.o   \
              /usr/local/lib/gcc/*/*/libgcc.a   /usr/local/lib/gcc/x86_64-pc-linux-musl/$VER/
COPY --from=0 /usr/local/lib64/libgfortran.a    \
              /usr/local/lib64/libgfortran.spec \
              /usr/local/lib64/libquadmath.a    /usr/local/lib64/
COPY --from=0 /usr/local/libexec/gcc/*/*/f951   /usr/local/libexec/gcc/x86_64-pc-linux-musl/

ENTRYPOINT ["fortran"]

CMD ["--version"]
