FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl

# Also rebuild Fennel.
RUN curl https://www.lua.org/ftp/lua-5.4.7.tar.gz | tar xz

RUN cd lua-*          \
 && make generic      \
    AR='gcc-ar rcu'   \
    MYCFLAGS=-flto    \
    MYLDFLAGS=-static \
    RANLIB=gcc-ranlib \
 && strip src/lua

FROM codegolf/lang-base-no-proc

COPY --from=0 /lua-*/src/lua /usr/bin/

ENTRYPOINT ["lua"]

CMD ["-v"]
