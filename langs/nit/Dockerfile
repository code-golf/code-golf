FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base ccache curl dash gc-dev libunwind-dev musl-dev pkgconfig

ENV VER=0.8

RUN curl -#L https://github.com/nitlang/nit/archive/refs/tags/v$VER.tar.gz | tar xz

RUN mv nit-$VER /usr/local/nit

WORKDIR /usr/local/nit

RUN CFLAGS='-Wno-cpp' make \
 && strip ./bin/nit

# Placeholder directory for the runner source file.
RUN mkdir -p /nit/src

COPY nitwrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/nitwrapper -s -static /nitwrapper.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/dash                           /bin/sh
COPY --from=0 /lib/ld-musl-x86_64.so.1                /lib/
COPY --from=0 /nit/src                                /nit/src
COPY --from=0 /usr/bin/as                             \
              /usr/bin/env                            \
              /usr/bin/gcc                            \
              /usr/bin/ld                             \
              /usr/bin/nitwrapper                     /usr/bin/
COPY --from=0 /usr/include                            /usr/include
COPY --from=0 /usr/lib/crti.o                         \
              /usr/lib/crtn.o                         \
              /usr/lib/libbfd-2.43.1.so               \
              /usr/lib/libc.so                        \
              /usr/lib/libctf.so.0                    \
              /usr/lib/libgc.so.1                     \
              /usr/lib/libgcc_s.so                    \
              /usr/lib/libgcc_s.so.1                  \
              /usr/lib/libgmp.so.10                   \
              /usr/lib/libisl.so.23                   \
              /usr/lib/libjansson.so.4                \
              /usr/lib/liblzma.so.5                   \
              /usr/lib/libm.a                         \
              /usr/lib/libmpc.so.3                    \
              /usr/lib/libmpfr.so.6                   \
              /usr/lib/libsframe.so.1                 \
              /usr/lib/libssp_nonshared.a             \
              /usr/lib/libunwind.so.8                 \
              /usr/lib/libz.so.1                      \
              /usr/lib/libzstd.so.1                   /usr/lib/
COPY --from=0 /usr/lib/gcc/x86_64*/14.2.0/crtbeginS.o \
              /usr/lib/gcc/x86_64*/14.2.0/crtendS.o   \
              /usr/lib/gcc/x86_64*/14.2.0/libgcc.a    /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/libexec/gcc                        /usr/libexec/gcc
COPY --from=0 /usr/local/nit                          /usr/local

ENTRYPOINT ["nitwrapper"]

CMD ["--version"]
