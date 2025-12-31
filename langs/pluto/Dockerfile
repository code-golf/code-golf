FROM alpine:3.23 AS builder

RUN apk add --no-cache curl g++ make

ENV VER=0.12.2

RUN curl -#L https://github.com/PlutoLang/Pluto/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

RUN make LDFLAGS='-static' PLAT='linux' \
 && make INSTALL_BIN='/usr/bin' install

RUN strip /usr/bin/pluto

FROM codegolf/lang-base

COPY --from=0 /usr/bin/pluto /usr/bin/

ENTRYPOINT ["pluto"]

CMD ["-v"]
