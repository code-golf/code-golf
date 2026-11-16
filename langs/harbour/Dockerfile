FROM alpine:3.24 AS builder

RUN apk add --no-cache gcc make musl-dev

ENV VER=3.2.0

WORKDIR /harbour

ADD https://github.com/harbour/core.git#v$VER .

RUN make -j`nproc` install \
    HB_BUILD_CONTRIBS='no' \
    HB_BUILD_STRIP='all'   \
    HB_INSTALL_PREFIX='/usr'

FROM codegolf/lang-base

COPY --from=0 --parents /lib/ld-musl-*.so.1  \
                        /usr/bin/hbmk2       \
                        /usr/include/harbour /

ENTRYPOINT ["hbmk2"]

CMD ["--version"]
