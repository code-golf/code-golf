FROM alpine:3.20 AS builder

ENV VER=8.20.0

RUN apk add --no-cache         \
             bash              \
             build-base        \
             cairo-dev         \
             curl              \
             git               \
             gmp-dev           \
             gtk+3.0-dev       \
             gtksourceview-dev \
             opam              \
             pkgconf

RUN opam init                                                                  \
   --disable-sandboxing                                                        \
   --shell-setup                                                               \
   --yes                                                                       \
 && eval $(opam env --switch=default)                                          \
 && opam switch create coq                                                     \
    --packages                                                                 \
    "ocaml-option-flambda","ocaml-variants.4.14.1+options"                     \
 && eval $(opam env --switch=coq)                                              \
 && opam install                                                               \
    --no-depexts                                                               \
    --yes                                                                      \
    "dune" "lablgtk3-sourceview3" "ocamlfind" "zarith"                         \
 && eval $(opam env)                                                           \
 && curl -L https://github.com/coq/coq/releases/download/V$VER/coq-$VER.tar.gz \
  | tar xz                                                                     \
 && mv coq-$VER coq                                                            \
 && cd coq                                                                     \
 && ./configure                                                                \
    -prefix "/usr"                                                             \
 && make dunestrap                                                             \
 && dune build                                                                 \
    -p                                                                         \
    "coq","coq-core","coq-stdlib"                                              \
 && dune install                                                               \
    --prefix "/usr"                                                            \
    "coq" "coq-core" "coq-stdlib"

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh                 /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1         /lib/
COPY --from=0 /root/.opam/coq/lib/findlib.conf /root/.opam/coq/lib/
COPY --from=0 /usr/bin/coqc                    /usr/bin/
COPY --from=0 /usr/lib/libgmp.so.10            /usr/lib/
COPY --from=0 /usr/lib/coq                     /usr/lib/coq
COPY --from=0 /usr/lib/coq-core                /usr/lib/coq-core
COPY --from=0 /usr/lib/coq-stdlib              /usr/lib/coq-stdlib

# coq:     82.6 MB
# coq-core: 556 MB

COPY /coqwrapper /usr/bin/

ENTRYPOINT ["coqwrapper"]

CMD ["--version"]
