ARG VER=10.3.0

FROM alpine:3.22 AS builder

ARG VER

RUN apk add --no-cache autoconf automake blas-dev curl g++ gawk gfortran gperf \
                       lapack-dev less libtool make ncurses-dev pcre2-dev

RUN curl -#L https://ftpmirror.gnu.org/gnu/octave/octave-$VER.tar.xz \
  | tar xJ --strip-components 1

RUN ./configure        \
    --disable-docs     \
    --disable-readline \
 && make -j`nproc` install

RUN strip /usr/local/bin/octave-cli-$VER

WORKDIR /usr/local/lib/octave/$VER

RUN strip liboctave.so.*    \
          liboctinterp.so.* \
          liboctmex.so.*

COPY octave.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/octave -s -static /octave.c

FROM codegolf/lang-base

ARG VER

COPY --from=0 /lib/ld-musl-*.so.1                           /lib/
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
COPY --from=0 /usr/local/lib/octave/$VER/liboctave.so.12    \
              /usr/local/lib/octave/$VER/liboctinterp.so.13 \
              /usr/local/lib/octave/$VER/liboctmex.so.1     /usr/local/lib/octave/$VER/
COPY --from=0 /usr/local/share/octave                       /usr/local/share/octave

ENTRYPOINT ["octave"]

CMD ["--quiet", "--version"]
