FROM alpine:3.23 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV VER=5.4.0

RUN curl -#L https://github.com/ocaml/ocaml/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure                           \
    --disable-debug-runtime               \
    --disable-installing-source-artifacts \
    --disable-ocamldebug                  \
    --disable-ocamldoc                    \
    --disable-ocamlobjinfo                \
    --disable-ocamltest                   \
    --disable-shared                      \
    --disable-stdlib-manpages             \
    --disable-systhreads                  \
    --enable-static                       \
    --prefix=/usr                         \
 && make -j`nproc` world.opt              \
 && make install                          \
 && strip /usr/bin/ocamlc                 \
          /usr/bin/ocamlrun

# Remove some libraries we don't need.
RUN rm -r /usr/lib/ocaml/compiler-libs

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-*.so.1 /lib/
COPY --from=0 /usr/bin/ocaml      \
              /usr/bin/ocamlc     \
              /usr/bin/ocamlrun   /usr/bin/
COPY --from=0 /usr/lib/ocaml      /usr/lib/ocaml

ENTRYPOINT ["ocaml"]

CMD ["--version"]
