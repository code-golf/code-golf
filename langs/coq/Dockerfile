FROM alpine:3.20 AS builder
#FROM ocaml/opam:alpine-3.20-ocaml-5.2-flambda AS builder

ENV VERSION="8.20.0"

RUN apk add --no-cache bash build-base cairo-dev git gmp-dev gtk+3.0-dev gtksourceview-dev opam pkgconf

RUN opam init                                           \
   --compiler="ocaml-base-compiler.5.2.0"               \
   --disable-sandboxing                                 \
   --jobs=`nproc`                                       \
   --shell-setup                                        \
   --yes                                                \
 && eval $(opam env --switch=ocaml-base-compiler.5.2.0) \
 && opam switch create coq                              \
    --jobs=`nproc`                                      \
    --packages                                          \
      "ocaml-option-flambda",                           \
      "ocaml-variants.4.14.1+options"                   \
 && eval $(opam env --switch=coq)                       \
 && opam install                                        \
    --jobs=`nproc`                                      \
    --no-depexts                                        \
    --yes                                               \
      "dune"                                            \
      "lablgtk3-sourceview3"                            \
      "ocamlfind"                                       \
      "zarith"                                          \
 && eval $(opam env)                                    \
 && git clone https://github.com/coq/coq.git            \
    --branch                                            \
       V$VERSION                                        \
 && cd coq                                              \
 && ./configure                                         \
    -prefix="/usr"                                      \
 && make dunestrap                                      \
 && dune build                                          \
    -p                                                  \
      "coq",                                            \
      "coq-core",                                       \
      "coq-stdlib"                                      \
 && dune install                                        \
    --prefix="/usr"                                     \
      "coq"                                             \
      "coq-core"                                        \
      "coq-stdlib"

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh                 /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1         /lib/
COPY --from=0 /root/.opam/coq/lib/findlib.conf /root/.opam/coq/lib/
COPY --from=0 /usr/bin/coqc                    /usr/bin/
COPY --from=0 /usr/lib/libacl.so.1             \
              /usr/lib/libattr.so.1            \
              /usr/lib/libcrypto.so.3          \
              /usr/lib/libgmp.so.10            \
              /usr/lib/libutmps.so.0.1         /usr/lib/
COPY --from=0 /usr/lib/coq                     /usr/lib/coq
COPY --from=0 /usr/lib/coq-core                /usr/lib/coq-core
COPY --from=0 /usr/lib/coq-stdlib              /usr/lib/coq-stdlib

# coq:     82.6 MB ...
# coq-core: 556 MB <<<<<<

COPY /coqwrapper /usr/bin/

ENTRYPOINT ["coqwrapper"]

CMD ["--version"]
