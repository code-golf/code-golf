FROM ocaml/opam:alpine-3.20-ocaml-5.2 AS builder

RUN sudo apk add gmp-dev linux-headers \
 && OPAMJOBS=`nproc` opam pin add coq 8.20.0

FROM codegolf/lang-base

COPY --from=0 /bin           /bin
COPY --from=0 /lib           /lib
COPY --from=0 /usr/bin       /usr/bin
COPY --from=0 /usr/lib       /usr/lib
COPY --from=0 /usr/local/bin /usr/local/bin
COPY --from=0 /usr/local/lib /usr/local/lib
COPY --from=0 /usr/share     /usr/share

COPY /coqwrapper /usr/bin/

ENTRYPOINT ["coqwrapper"]

CMD ["--version"]

#COPY --from=0 /bin/cat /bin/sh                 /bin/
#COPY --from=0 /lib/ld-musl-x86_64.so.1         /lib/
#COPY --from=0 /root/.opam/coq/lib/findlib.conf /root/.opam/coq/lib/
#COPY --from=0 /usr/bin/coqc                    /usr/bin/
#COPY --from=0 /usr/lib/libacl.so.1             \
#              /usr/lib/libattr.so.1            \
#              /usr/lib/libcrypto.so.3          \
#              /usr/lib/libgmp.so.10            \
#              /usr/lib/libutmps.so.0.1         /usr/lib/
#COPY --from=0 /usr/lib/coq                     /usr/lib/coq
#COPY --from=0 /usr/lib/coq-core                /usr/lib/coq-core
#COPY --from=0 /usr/lib/coq-stdlib              /usr/lib/coq-stdlib
