FROM alpine:edge

RUN apk add --no-cache curl gcc make musl-dev

RUN curl https://www.lua.org/ftp/lua-5.3.4.tar.gz | tar xzf -

RUN cd lua-5.3.4                   \
 && make generic MYLDFLAGS=-static \
 && strip src/lua

FROM scratch

COPY --from=0 /lua-5.3.4/src/lua /usr/bin/

ENTRYPOINT ["/usr/bin/lua", "-v"]
