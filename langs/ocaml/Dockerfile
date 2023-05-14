FROM alpine:3.18 AS builder

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/ocaml/ocaml/archive/refs/tags/5.0.0.tar.gz | tar xz

RUN cd ocaml-5.0.0          \
 && ./configure             \
    --bindir /usr/bin       \
    --libdir /usr/lib/ocaml \
    --mandir /usr/share/man \
    --prefix /usr           \
 && make world.opt          \
 && make install            \
 && strip /usr/bin/ocamlc /usr/bin/ocamlrun

# Remove some crap we don't need.
RUN rm -rf /usr/lib/ocaml/compiler-libs /usr/lib/ocaml/ocamldoc \
 && find /usr/lib/ocaml \
    \( -name '*.cmt' -o -name '*.cmti' -o -name '*.ml' \) -a -delete

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/ocaml           \
              /usr/bin/ocamlc          \
              /usr/bin/ocamlrun        /usr/bin/
COPY --from=0 /usr/lib/ocaml           /usr/lib/ocaml

ENTRYPOINT ["ocaml"]

CMD ["--version"]
