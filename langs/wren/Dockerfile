FROM alpine:3.22 AS builder

RUN apk add --no-cache build-base curl

ENV LDFLAGS='-static' VER=0.4.0

RUN curl -#L https://github.com/wren-lang/wren-cli/archive/refs/tags/$VER.tar.gz \
  | tar xz --strip-components 1

COPY stdin.patch /

RUN patch -p0 < stdin.patch

WORKDIR /projects/make

RUN make

FROM codegolf/lang-base

COPY --from=0 /bin/wren_cli /usr/bin/wren

ENTRYPOINT ["wren"]

CMD ["--version"]
