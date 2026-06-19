FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev

ENV VER=5.5.0

ADD --unpack https://github.com/ocaml/ocaml/archive/refs/tags/$VER.tar.gz /

WORKDIR /ocaml-$VER

RUN ./configure               \
    --disable-debug-runtime   \
    --disable-ocamldebug      \
    --disable-ocamldoc        \
    --disable-ocamlobjinfo    \
    --disable-ocamltest       \
    --disable-stdlib-manpages \
    --disable-systhreads      \
    --enable-static           \
    --prefix=/usr             \
 && make -j`nproc` world.opt  \
 && make install              \
 && strip /usr/bin/ocamlc     \
          /usr/bin/ocamlrun

WORKDIR /usr/lib/ocaml

# Remove some libraries and source artifacts we don't need.
RUN rm -r compiler-libs \
 && find . \( -name '*.cmt*' -o -name '*.ml' \) -a -delete

RUN rm /usr/bin/ocamldep* /usr/bin/ocamlopt*

FROM codegolf/lang-base

COPY --from=0           /lib/ld-musl-*.so.1             /lib/
COPY --from=0 --parents /usr/bin/*ocaml* /usr/lib/ocaml /

ENTRYPOINT ["ocaml"]

CMD ["--version"]
