FROM alpine:3.23 AS builder

RUN apk add --no-cache curl gcc make musl-dev patch

ENV LDFLAGS='-static' VER=0.4.0

RUN curl -#L https://github.com/wren-lang/wren-cli/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

COPY wren.patch /

RUN patch -p0 < wren.patch

WORKDIR /projects/make

RUN make

FROM codegolf/lang-base

COPY --from=0 /bin/wren_cli /usr/bin/wren

ENTRYPOINT ["wren"]

CMD ["--version"]
