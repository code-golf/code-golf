FROM alpine:3.12 as builder

RUN apk add --no-cache curl g++ libexecinfo-dev \
    linux-headers m4 make patch tar xz

WORKDIR /src

RUN curl -L ftp://ftp.gnu.org/gnu/gcc/gcc-9.3.0/gcc-9.3.0.tar.gz | tar xzf -
RUN curl -L ftp://gcc.gnu.org/pub/gcc/infrastructure/gmp-4.3.2.tar.bz2 | tar xjf -
RUN curl -L ftp://gcc.gnu.org/pub/gcc/infrastructure/mpfr-2.4.2.tar.bz2 | tar xjf -
RUN curl -L ftp://gcc.gnu.org/pub/gcc/infrastructure/mpc-0.8.1.tar.gz | tar xzf -

WORKDIR /src/gmp-4.3.2

RUN ./configure --disable-shared --enable-static --prefix=/tmp/gcc
RUN make
RUN make install

WORKDIR /src/mpfr-2.4.2

RUN ./configure --disable-shared --enable-static --prefix=/tmp/gcc --with-gmp=/tmp/gcc
RUN make
RUN make install

WORKDIR /src/mpc-0.8.1

RUN ./configure --disable-shared --enable-static --prefix=/tmp/gcc --with-gmp=/tmp/gcc --with-mpfr=/tmp/gcc
RUN make
RUN make install

WORKDIR /scratch

RUN /src/gcc-9.3.0/configure   \
    --disable-bootstrap        \
    --disable-libstdcxx-pch    \
    --disable-multilib         \
    --disable-shared           \
    --enable-languages=fortran \
    --enable-libgomp           \
    --enable-lto               \
    --enable-static            \
    --enable-threads=posix     \
    --enable-tls               \
    --prefix=/fortran          \
    --with-fpmath=sse          \
    --with-gmp=/tmp/gcc        \
    --with-mpc=/tmp/gcc        \
    --with-mpfr=/tmp/gcc       \
 && make                       \
 && make install

FROM scratch

COPY --from=0 /fortran/bin/            /fortran/bin/
COPY --from=0 /fortran/libexec/gcc/x86_64-pc-linux-musl/9.3.0/f951                   \
              /fortran/libexec/gcc/x86_64-pc-linux-musl/9.3.0/liblto_plugin.la       \
              /fortran/libexec/gcc/x86_64-pc-linux-musl/9.3.0/liblto_plugin.so       \
              /fortran/libexec/gcc/x86_64-pc-linux-musl/9.3.0/liblto_plugin.so.0     \
              /fortran/libexec/gcc/x86_64-pc-linux-musl/9.3.0/liblto_plugin.so.0.0.0 /fortran/libexec/gcc/x86_64-pc-linux-musl/9.3.0/
COPY --from=0 /fortran/lib/            /fortran/lib/
COPY --from=0 /fortran/lib64/          /fortran/lib64/
COPY --from=0 /usr/bin/                /usr/bin/
COPY --from=0 /usr/lib/                /usr/lib/
COPY --from=0 /bin/cat                 \
              /bin/sh                  \
              /bin/rm                  /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1 \
              /lib/libz.so.1           /lib/

COPY fortran                           /usr/bin/

ENTRYPOINT ["/fortran/bin/gfortran", "--version"]
