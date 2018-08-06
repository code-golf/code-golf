FROM alpine:edge as builder

RUN apk add --no-cache curl gcc make musl-dev

RUN curl https://www.lua.org/ftp/lua-5.3.5.tar.gz | tar xzf -

RUN cd lua-5.3.5      \
 && make generic      \
    AR='gcc-ar rcu'   \
    MYCFLAGS=-flto    \
    MYLDFLAGS=-static \
    RANLIB=gcc-ranlib \
 && strip src/lua

FROM scratch

COPY --from=0 /lua-5.3.5/src/lua /usr/bin/

ENTRYPOINT ["/usr/bin/lua", "-v"]
