FROM alpine:3.24 AS builder

RUN apk add --no-cache bash gcc git make musl-dev

ENV VER=0.12.0

WORKDIR /CBQN

ADD --keep-git-dir https://github.com/dzaima/CBQN.git#v$VER .

RUN make          \
    REPLXX=0      \
    singeli=0     \
    static-bin    \
    version=v$VER \
 && make install  \
 && strip /usr/local/bin/bqn

FROM codegolf/lang-base

COPY --from=0 /usr/local/bin/bqn /usr/bin/

ENTRYPOINT ["bqn"]

CMD ["--version"]
