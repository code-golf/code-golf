FROM alpine:3.15 as builder

RUN mkdir /empty

RUN apk add --no-cache build-base curl

RUN curl https://www.lua.org/ftp/lua-5.4.4.tar.gz | tar xz

RUN cd lua-5.4.4      \
 && make generic      \
    AR='gcc-ar rcu'   \
    MYCFLAGS=-flto    \
    MYLDFLAGS=-static \
    RANLIB=gcc-ranlib \
 && strip src/lua

FROM scratch

COPY --from=0 /empty             /dev
COPY --from=0 /empty             /proc
COPY --from=0 /empty             /tmp
COPY --from=0 /lua-5.4.4/src/lua /usr/bin/

ENTRYPOINT ["lua"]

CMD ["-v"]
