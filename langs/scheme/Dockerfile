FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev

ENV VER=10.4.1

ADD https://github.com/cisco/ChezScheme.git#v$VER /

RUN ./configure      \
    --disable-curses \
    --disable-x11    \
 && make -j`nproc`   \
 && make install

RUN strip /usr/bin/scheme

FROM codegolf/lang-base

COPY --from=0           /usr/bin/scheme     /usr/bin/
COPY --from=0 --parents /boot/ta6le/*.boot  \
                        /lib/ld-musl-*.so.1 /

ENTRYPOINT ["scheme"]

CMD ["--version"]
