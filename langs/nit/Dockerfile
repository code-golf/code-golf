FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base ccache curl gc-dev libunwind-dev musl-dev pkgconfig

RUN curl -L https://github.com/nitlang/nit/archive/refs/tags/v0.8.tar.gz \
  | tar xz

RUN mv nit-* /usr/local/nit \
 && cd /usr/local/nit       \
 && CFLAGS=-Wno-cpp make

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/mkdir /bin/sh  /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1     \
              /lib/libz.so.1               /lib/
COPY --from=0 /usr/bin/as                  \
              /usr/bin/env                 \
              /usr/bin/gcc                 \
              /usr/bin/ld                  /usr/bin/
COPY --from=0 /usr/include/alloca.h        \
              /usr/include/dirent.h        \
              /usr/include/endian.h        \
              /usr/include/errno.h         \
              /usr/include/features.h      \
              /usr/include/inttypes.h      \
              /usr/include/math.h          \
              /usr/include/poll.h          \
              /usr/include/stdint.h        \
              /usr/include/stdio.h         \
              /usr/include/stdlib.h        \
              /usr/include/string.h        \
              /usr/include/strings.h       \
              /usr/include/time.h          \
              /usr/include/unistd.h        /usr/include/
COPY --from=0 /usr/include/bits/alltypes.h \
              /usr/include/bits/dirent.h   \
              /usr/include/bits/errno.h    \
              /usr/include/bits/poll.h     \
              /usr/include/bits/posix.h    \
              /usr/include/bits/stat.h     \
              /usr/include/bits/stdint.h   /usr/include/bits/
COPY --from=0 /usr/include/sys/select.h    \
              /usr/include/sys/stat.h      \
              /usr/include/sys/types.h     /usr/include/sys/
COPY --from=0 /usr/lib/crti.o              \
              /usr/lib/crtn.o              \
              /usr/lib/libbfd-2.42.so      \
              /usr/lib/libc.so             \
              /usr/lib/libctf.so.0         \
              /usr/lib/libgc.so.1          \
              /usr/lib/libgcc_s.so         \
              /usr/lib/libgcc_s.so.1       \
              /usr/lib/libgmp.so.10        \
              /usr/lib/libisl.so.23        \
              /usr/lib/libjansson.so.4     \
              /usr/lib/liblzma.so.5        \
              /usr/lib/libm.a              \
              /usr/lib/libmpc.so.3         \
              /usr/lib/libmpfr.so.6        \
              /usr/lib/libsframe.so.1      \
              /usr/lib/libssp_nonshared.a  \
              /usr/lib/libunwind.so.8      \
              /usr/lib/libzstd.so.1        \
              /usr/lib/gcc/*/*/crtbeginS.o \
              /usr/lib/gcc/*/*/crtendS.o   \
              /usr/lib/gcc/*/*/libgcc.a    /usr/lib/
COPY --from=0 /usr/libexec/gcc             /usr/lib/gcc
COPY --from=0 /usr/local/nit               /usr/local

COPY /nitwrapper /usr/bin/

ENTRYPOINT ["nitwrapper"]

CMD ["--version"]
