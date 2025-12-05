ARG MLVER=5.4 OPAMHOME=/home/opam/.opam/$MLVER

FROM ocaml/opam:alpine-3.22-ocaml-$MLVER AS builder

RUN sudo apk add --no-cache gmp-dev linux-headers

ENV VER=9.1.0

RUN opam pin -j`nproc` add rocq-core $VER

COPY rocq.c /

RUN sudo gcc -Wall -Werror -Wextra -o /usr/bin/rocq -s -static /rocq.c

FROM codegolf/lang-base

ARG OPAMHOME

COPY --from=0 /lib/ld-musl-x86_64.so.1    /lib/
COPY --from=0 /usr/bin/rocq               /usr/bin/
COPY --from=0 /usr/lib/libacl.so.1        \
              /usr/lib/libattr.so.1       \
              /usr/lib/libcrypto.so.3     \
              /usr/lib/libgmp.so.10       \
              /usr/lib/libskarnet.so.2.14 \
              /usr/lib/libutmps.so.0.1    /usr/lib/
COPY --from=0 $OPAMHOME/bin/rocq          /usr/local/bin/
COPY --from=0 $OPAMHOME/lib/findlib.conf  /usr/local/lib/
COPY --from=0 $OPAMHOME/lib/coq/theories  /usr/local/lib/coq/theories

ENTRYPOINT ["rocq"]

CMD ["--version"]
