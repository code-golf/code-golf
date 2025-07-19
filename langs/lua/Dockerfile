FROM alpine:3.22 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV VER=5.4.8

# Also rebuild Fennel.
RUN curl -#L https://www.lua.org/ftp/lua-$VER.tar.gz \
  | tar xz --strip-components 1

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
