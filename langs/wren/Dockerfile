FROM alpine:3.16 as builder

RUN apk add --no-cache build-base curl

RUN curl -L https://github.com/wren-lang/wren-cli/archive/refs/tags/0.4.0.tar.gz \
  | tar xz

RUN cd wren-cli-0.4.0/projects/make && make

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh             /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1     /lib/
COPY --from=0 /wren-cli-0.4.0/bin/wren_cli /usr/bin/

COPY wren /usr/bin/

ENTRYPOINT ["wren"]

CMD ["--version"]
