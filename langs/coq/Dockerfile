FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl gmp-dev libgtksourceview-3.0-dev opam

ENV VER=8.20.0

RUN opam init -y --disable-sandboxing --shell-setup                                        \
 && eval $(opam env --switch=default)                                                      \
 && opam switch create coq --packages="ocaml-variants.4.14.1+options,ocaml-option-flambda" \
 && eval $(opam env --switch=coq)                                                          \
 && opam install -y dune ocamlfind zarith lablgtk3-sourceview3                             \
 && eval $(opam env)                                                                       \
 && curl -L https://github.com/coq/coq/releases/download/V$VER/coq-$VER.tar.gz | tar xz    \
 && mv coq-$VER coq                                                                        \
 && cd coq                                                                                 \
 && ./configure                                                                            \
    -prefix /usr                                                                           \
 && make dunestrap                                                                         \
 && dune build -p coq-core,coq-stdlib,coq,coqide-server,coqide                             \
 && dune install --prefix=/usr coq-core coq-stdlib coq coqide-server

FROM codegolf/lang-base

COPY --from=0 / /

ENTRYPOINT ["/usr/lib/coq"]

CMD ["--version"]
