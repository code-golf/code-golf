FROM alpine:3.12 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base curl db-dev g++ gmp-dev gnupg \
    make ncurses-dev patch tar xz

RUN curl https://ftp.gnu.org/gnu/gnucobol/gnucobol-2.2.tar.xz -O
RUN curl https://ftp.gnu.org/gnu/gnucobol/gnucobol-2.2.tar.xz.sig -O
RUN curl https://ftp.gnu.org/gnu/gnu-keyring.gpg -O
RUN gpg --verify --keyring ./gnu-keyring.gpg gnucobol-2.2.tar.xz.sig
RUN tar xJf gnucobol-2.2.tar.xz

WORKDIR /gnucobol-2.2
RUN ./configure
RUN make
RUN make install
RUN make check

WORKDIR /
RUN curl -L https://download.savannah.nongnu.org/releases/tinycc/tcc-0.9.27.tar.bz2 | tar xjf -

COPY setvbuf.patch /

RUN mv tcc-* tcc             \
 && patch -p0 <setvbuf.patch \
 && cd tcc                   \
 && ./configure              \
    --config-musl            \
    --prefix=/usr            \
    --strip-binaries         \
 && make                     \
 && make install

COPY cobol.c /
RUN gcc -s -o cobol cobol.c

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1    /lib/
COPY --from=0 /empty                      /proc
COPY --from=0 /empty                      /tmp
COPY --from=0 /usr/local/bin/cobc         /usr/local/bin/
COPY --from=0 /usr/local/share/gnucobol/  /usr/local/share/gnucobol/
COPY --from=0 /usr/local/include/libcob.h /usr/local/include/
COPY --from=0 /usr/local/include/libcob/  /usr/local/include/libcob/
COPY --from=0 /usr/local/lib/libcob.so    \
              /usr/local/lib/libcob.so.4  /usr/local/lib/
COPY --from=0 /usr/bin/tcc                \
              /cobol                      /usr/bin/
COPY --from=0 /usr/lib/libdb-5.3.so       \
              /usr/lib/libgmp.so.10       \
              /usr/lib/libncursesw.so.6   /usr/lib/
COPY --from=0 /usr/include/assert.h       \
              /usr/include/alloca.h       \
              /usr/include/complex.h      \
              /usr/include/ctype.h        \
              /usr/include/errno.h        \
              /usr/include/features.h     \
              /usr/include/fenv.h         \
              /usr/include/float.h        \
              /usr/include/gmp.h          \
              /usr/include/inttypes.h     \
              /usr/include/iso646.h       \
              /usr/include/limits.h       \
              /usr/include/locale.h       \
              /usr/include/math.h         \
              /usr/include/setjmp.h       \
              /usr/include/signal.h       \
              /usr/include/stdalign.h     \
              /usr/include/stdarg.h       \
              /usr/include/stdbool.h      \
              /usr/include/stddef.h       \
              /usr/include/stdint.h       \
              /usr/include/stdio.h        \
              /usr/include/stdlib.h       \
              /usr/include/stdnoreturn.h  \
              /usr/include/string.h       \
              /usr/include/strings.h      \
              /usr/include/tgmath.h       \
              /usr/include/threads.h      \
              /usr/include/time.h         \
              /usr/include/uchar.h        \
              /usr/include/wchar.h        \
              /usr/include/wctype.h       /usr/include/
COPY --from=0 /usr/include/bits           /usr/include/bits/
COPY --from=0 /usr/lib/libc.so            /usr/lib/
COPY --from=0 /usr/lib/tcc/libtcc1.a      /usr/lib/tcc/

ENTRYPOINT ["/usr/local/bin/cobc", "--version"]
