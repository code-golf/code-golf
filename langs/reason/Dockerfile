FROM codegolf/lang-ocaml

FROM alpine:3.22 AS builder

RUN apk add --no-cache bash dash gcc make musl-dev patch

ENV BIN='/usr/bin/opam' OPAMVER=2.4.1 VER=3.17.0

RUN wget -qO $BIN https://github.com/ocaml/opam/releases/download/$OPAMVER/opam-$OPAMVER-x86_64-linux \
 && chmod +x $BIN

RUN $BIN init --disable-sandboxing --yes \
 && $BIN pin -j`nproc` add reason $VER --yes

COPY reason.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/reason -s reason.c

FROM codegolf/lang-base

COPY --from=0 /usr/lib/ocaml                        /root/.opam/default/lib/ocaml
COPY --from=1 /usr/bin/dash                         /bin/sh
COPY --from=1 /lib/ld-musl-*.so.1                   /lib/
COPY --from=1 /usr/bin/as                           \
              /usr/bin/gcc                          \
              /usr/bin/ld                           \
              /usr/bin/reason                       /usr/bin/
COPY --from=1 /usr/include                          /usr/include
COPY --from=1 /usr/lib/Scrt1.o                      \
              /usr/lib/crti.o                       \
              /usr/lib/crtn.o                       \
              /usr/lib/libbfd-2.44.so               \
              /usr/lib/libc.so                      \
              /usr/lib/libctf.so.0                  \
              /usr/lib/libgcc_s.so                  \
              /usr/lib/libgcc_s.so.1                \
              /usr/lib/libgmp.so.10                 \
              /usr/lib/libisl.so.23                 \
              /usr/lib/libjansson.so.4              \
              /usr/lib/libmpc.so.3                  \
              /usr/lib/libmpfr.so.6                 \
              /usr/lib/libpthread.a                 \
              /usr/lib/libsframe.so.1               \
              /usr/lib/libssp_nonshared.a           \
              /usr/lib/libz.so.1                    \
              /usr/lib/libzstd.so.1                 /usr/lib/
COPY --from=1 /usr/lib/gcc/*/*/crtbeginS.o          \
              /usr/lib/gcc/*/*/crtendS.o            \
              /usr/lib/gcc/*/*/libgcc.a             /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=1 /usr/lib/gcc/*/*/include/stdatomic.h  /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/include/
COPY --from=1 /usr/libexec/gcc/*/*/cc1              \
              /usr/libexec/gcc/*/*/liblto_plugin.so /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=1 /root/.opam/default/bin/dune          \
              /root/.opam/default/bin/ocamlc        \
              /root/.opam/default/bin/refmt         /usr/local/bin/

ENTRYPOINT ["reason"]
