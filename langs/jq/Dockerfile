FROM alpine:3.22 AS builder

RUN apk add --no-cache autoconf automake build-base git libtool

# Three commits past 1.8.0 to fix https://github.com/jqlang/jq/issues/3339
ENV VER=023f274

RUN git clone https://github.com/jqlang/jq.git

WORKDIR /jq

RUN git reset --hard $VER              \
 && git submodule update --init        \
 && autoreconf -i                      \
 && ./configure                        \
    --prefix=/usr                      \
    --with-oniguruma=builtin           \
 && make LDFLAGS='-all-static' install \
 && strip /usr/bin/jq

FROM codegolf/lang-base

COPY --from=0 /usr/bin/jq /usr/bin/

ENTRYPOINT ["jq"]

CMD ["--version"]
