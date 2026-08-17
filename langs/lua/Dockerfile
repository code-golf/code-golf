FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev

# Also rebuild Fennel.
ENV VER=5.5.1

ADD --unpack https://www.lua.org/ftp/lua-$VER.tar.gz /

WORKDIR /lua-$VER

RUN make generic install \
    AR='gcc-ar rcu'      \
    MYCFLAGS='-flto'     \
    MYLDFLAGS='-static'  \
    RANLIB='gcc-ranlib'  \
 && strip /usr/local/bin/lua

FROM codegolf/lang-base

COPY --from=0 /usr/local/bin/lua /usr/bin/

ENTRYPOINT ["lua"]

CMD ["-v"]
