FROM alpine:3.16 as builder

RUN apk add --no-cache build-base curl

RUN curl https://www.lua.org/ftp/lua-5.4.4.tar.gz | tar xz

RUN cd lua-5.4.4      \
 && make generic      \
    AR='gcc-ar rcu'   \
    MYCFLAGS=-flto    \
    MYLDFLAGS=-static \
    RANLIB=gcc-ranlib \
 && strip src/lua

FROM codegolf/lang-base

COPY --from=0 /lua-5.4.4/src/lua /usr/bin/

ENTRYPOINT ["lua"]

CMD ["-v"]
