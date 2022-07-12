FROM alpine:3.16 as builder

RUN apk add --no-cache build-base curl

RUN curl -L https://download.savannah.nongnu.org/releases/tinycc/tcc-0.9.27.tar.bz2 \
  | tar xj

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

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1   /lib/
COPY --from=0 /usr/bin/tcc               /usr/bin/
COPY --from=0 /usr/include/alloca.h      \
              /usr/include/assert.h      \
              /usr/include/complex.h     \
              /usr/include/ctype.h       \
              /usr/include/errno.h       \
              /usr/include/features.h    \
              /usr/include/fenv.h        \
              /usr/include/float.h       \
              /usr/include/inttypes.h    \
              /usr/include/iso646.h      \
              /usr/include/limits.h      \
              /usr/include/locale.h      \
              /usr/include/math.h        \
              /usr/include/setjmp.h      \
              /usr/include/signal.h      \
              /usr/include/stdalign.h    \
              /usr/include/stdarg.h      \
              /usr/include/stdbool.h     \
              /usr/include/stddef.h      \
              /usr/include/stdint.h      \
              /usr/include/stdio.h       \
              /usr/include/stdlib.h      \
              /usr/include/stdnoreturn.h \
              /usr/include/string.h      \
              /usr/include/strings.h     \
              /usr/include/tgmath.h      \
              /usr/include/threads.h     \
              /usr/include/time.h        \
              /usr/include/uchar.h       \
              /usr/include/wchar.h       \
              /usr/include/wctype.h      /usr/include/
COPY --from=0 /usr/include/bits          /usr/include/bits/
COPY --from=0 /usr/lib/libc.so           /usr/lib/
COPY --from=0 /usr/lib/tcc/libtcc1.a     /usr/lib/tcc/

ENTRYPOINT ["tcc"]

CMD ["-v"]
