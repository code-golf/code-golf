FROM ocaml/opam:alpine-3.20-ocaml-5.2 AS builder

RUN sudo apk add gmp-dev linux-headers         \
 && opam install -j`nproc` coq-stdlib="8.20.0" \
 && strip /home/opam/.opam/5.2/bin/coqc

FROM codegolf/lang-base

ENV OPAM_HOME="/home/opam/.opam/5.2"

COPY --from=0 /bin/cat /bin/sh                /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1        \
              /lib/libacl.so.1                \
              /lib/libattr.so.1               \
              /lib/libcrypto.so.3             \
              /lib/libskarnet.so.2.14         \
              /lib/libutmps.so.0.1            /lib/
COPY --from=0 $OPAM_HOME/bin/coqc             /usr/bin/
COPY --from=0 /usr/lib/libgmp.so.10           /usr/lib/
COPY --from=0 $OPAM_HOME/lib/findlib.conf     $OPAM_HOME/lib/
COPY --from=0 $OPAM_HOME/lib/coq/theories     $OPAM_HOME/lib/coq/theories
COPY --from=0 $OPAM_HOME/lib/coq-core/plugins $OPAM_HOME/lib/coq-core/plugins
COPY --from=0 $OPAM_HOME/lib/coq-stdlib/opam  $OPAM_HOME/lib/coq-stdlib/

COPY /coqwrapper /usr/bin/

ENTRYPOINT ["coqwrapper"]

CMD ["--version"]
