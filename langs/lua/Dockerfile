FROM alpine:3.19 AS builder

RUN apk add --no-cache build-base curl

RUN curl https://www.lua.org/ftp/lua-5.4.6.tar.gz | tar xz

RUN cd lua-5.4.6      \
 && make generic      \
    AR='gcc-ar rcu'   \
    MYCFLAGS=-flto    \
    MYLDFLAGS=-static \
    RANLIB=gcc-ranlib \
 && strip src/lua

FROM codegolf/lang-base-no-proc

COPY --from=0 /lua-5.4.6/src/lua /usr/bin/

ENTRYPOINT ["lua"]

CMD ["-v"]
