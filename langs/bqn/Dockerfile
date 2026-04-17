FROM alpine:3.23 AS builder

RUN apk add --no-cache bash clang git make

ENV VER=0.11.0

RUN git clone --branch v$VER --progress https://github.com/dzaima/CBQN

WORKDIR /CBQN

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
