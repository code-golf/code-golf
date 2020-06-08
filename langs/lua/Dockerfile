FROM alpine:3.11 as builder

RUN apk add --no-cache build-base curl

RUN curl https://www.lua.org/work/lua-5.4.0-rc4.tar.gz | tar xzf -

RUN cd lua-5.4.0      \
 && make generic      \
    AR='gcc-ar rcu'   \
    MYCFLAGS=-flto    \
    MYLDFLAGS=-static \
    RANLIB=gcc-ranlib \
 && strip src/lua

FROM scratch

COPY --from=0 /lua-5.4.0/src/lua /usr/bin/

ENTRYPOINT ["/usr/bin/lua", "-v"]
