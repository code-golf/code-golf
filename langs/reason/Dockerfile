ARG MLVER=5.3 OPAMHOME=/home/opam/.opam/$MLVER

FROM ocaml/opam:alpine-3.21-ocaml-$MLVER AS builder

RUN sudo apk add --no-cache dash

ENV VER=3.14.0

RUN OPAMJOBS=`nproc` opam pin add reason $VER

WORKDIR /reason

RUN echo '(executable (name code))' > dune \
 && echo '(lang dune 3.17)' > dune-project

COPY reason.c /

RUN sudo gcc -Wall -Werror -Wextra -o /usr/bin/reason -s /reason.c

FROM codegolf/lang-base

ARG OPAMHOME

COPY --from=0 /usr/bin/dash                              /bin/sh
COPY --from=0 /lib/ld-musl-x86_64.so.1                   /lib/
COPY --from=0 /reason/dune /reason/dune-project          /reason/
COPY --from=0 /usr/bin/as                                \
              /usr/bin/gcc                               \
              /usr/bin/ld                                \
              /usr/bin/reason                            /usr/bin/
COPY --from=0 /usr/include                               /usr/include
COPY --from=0 /usr/lib/Scrt1.o                           \
              /usr/lib/crti.o                            \
              /usr/lib/crtn.o                            \
              /usr/lib/libbfd-2.43.1.so                  \
              /usr/lib/libc.so                           \
              /usr/lib/libctf.so.0                       \
              /usr/lib/libgcc_s.so                       \
              /usr/lib/libgcc_s.so.1                     \
              /usr/lib/libgmp.so.10                      \
              /usr/lib/libisl.so.23                      \
              /usr/lib/libjansson.so.4                   \
              /usr/lib/libmpc.so.3                       \
              /usr/lib/libmpfr.so.6                      \
              /usr/lib/libpthread.a                      \
              /usr/lib/libsframe.so.1                    \
              /usr/lib/libssp_nonshared.a                \
              /usr/lib/libz.so.1                         \
              /usr/lib/libzstd.so.1                      /usr/lib/
COPY --from=0 /usr/lib/gcc/*/14.2.0/crtbeginS.o          \
              /usr/lib/gcc/*/14.2.0/crtendS.o            \
              /usr/lib/gcc/*/14.2.0/libgcc.a             /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 /usr/lib/gcc/*/14.2.0/include/stdatomic.h  /usr/lib/gcc/x86_64-alpine-linux-musl/14.2.0/include/
COPY --from=0 /usr/libexec/gcc/*/14.2.0/cc1              \
              /usr/libexec/gcc/*/14.2.0/liblto_plugin.so /usr/libexec/gcc/x86_64-alpine-linux-musl/14.2.0/
COPY --from=0 $OPAMHOME/bin/dune                         \
              $OPAMHOME/bin/ocamlc                       \
              $OPAMHOME/bin/refmt                        /usr/local/bin/
COPY --from=0 $OPAMHOME/lib/ocaml                        $OPAMHOME/lib/ocaml

ENTRYPOINT ["reason"]

CMD ["--version"]
