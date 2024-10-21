FROM alpine:3.20 AS builder

RUN apk add --no-cache bash build-base clang git

ENV VER=0.7.0

RUN git clone https://github.com/dzaima/CBQN.git --branch v$VER

RUN cd CBQN       \
 && make          \
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
