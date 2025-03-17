ARG VER=9.4.0

FROM alpine:3.21 AS builder

ARG VER

RUN apk add --no-cache autoconf automake blas-dev build-base curl gawk gfortran \
                       gperf lapack-dev less libtool ncurses pcre2 pcre2-dev

RUN curl -#L https://ftp.nluug.nl/pub/gnu/octave/octave-$VER.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure        \
    --disable-docs     \
    --disable-readline \
 && make -j`nproc` install

WORKDIR /usr/local/lib/octave/$VER

RUN strip liboctave.so.11 \
          liboctinterp.so.12

COPY octave.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/octave -s -static /octave.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/ld-musl-x86_64.so.1                      /lib/
COPY --from=0 /usr/bin/octave                               /usr/bin/
COPY --from=0 /usr/lib/libblas.so.3                         \
              /usr/lib/libgcc_s.so.1                        \
              /usr/lib/libgfortran.so.5                     \
              /usr/lib/libgomp.so.1                         \
              /usr/lib/liblapack.so.3                       \
              /usr/lib/libpcre2-8.so.0                      \
              /usr/lib/libquadmath.so.0                     \
              /usr/lib/libstdc++.so.6                       \
              /usr/lib/libz.so.1                            /usr/lib/
COPY --from=0 /usr/local/bin/octave                         \
              /usr/local/bin/octave-cli-$VER                /usr/local/bin/
COPY --from=0 /usr/local/lib/octave/$VER/liboctave.so.11    \
              /usr/local/lib/octave/$VER/liboctinterp.so.12 /usr/local/lib/octave/$VER/
COPY --from=0 /usr/local/share/octave                       /usr/local/share/octave

ENTRYPOINT ["octave"]

CMD ["--version"]
