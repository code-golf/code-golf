FROM alpine:3.18 AS builder

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/wren-lang/wren-cli/archive/refs/tags/0.4.0.tar.gz \
  | tar xz

COPY stdin.patch /

RUN mv wren-* wren         \
 && patch -p0 <stdin.patch \
 && LDFLAGS=-static make -C /wren/projects/make

FROM codegolf/lang-base-no-proc

COPY --from=0 /wren/bin/wren_cli /usr/bin/wren

ENTRYPOINT ["wren"]

CMD ["--version"]
