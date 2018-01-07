FROM alpine:edge

RUN apk add --no-cache curl gcc make musl-dev readline-dev

RUN curl https://www.lua.org/ftp/lua-5.3.4.tar.gz | tar xzf -

RUN cd lua-5.3.4                                \
 && sed -i -e /LUA_USE_READLINE/d src/luaconf.h \
 && make linux

FROM scratch

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /lua-5.3.4/src/lua       /usr/bin/

ENTRYPOINT ["/usr/bin/lua", "-v"]
