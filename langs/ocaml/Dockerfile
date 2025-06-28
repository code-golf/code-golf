FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl

ENV VER=5.3.0

RUN curl -#L https://github.com/ocaml/ocaml/archive/refs/tags/$VER.tar.gz | tar xz

WORKDIR /ocaml-$VER

RUN ./configure              \
    --prefix=/usr            \
 && make -j`nproc` world.opt \
 && make install             \
 && strip /usr/bin/ocamlc    \
          /usr/bin/ocamlrun

WORKDIR /usr/lib/ocaml

# Remove some stuff we don't need.
RUN rm -rf compiler-libs ocamldoc \
 && find . \( -name '*.cmt*' -o -name '*.ml' \) -a -delete

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/ocaml           \
              /usr/bin/ocamlc          \
              /usr/bin/ocamlrun        /usr/bin/
COPY --from=0 /usr/lib/ocaml           /usr/lib/ocaml

ENTRYPOINT ["ocaml"]

CMD ["--version"]
