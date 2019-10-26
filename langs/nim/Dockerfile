FROM alpine:3.10 as builder

RUN apk add -X http://dl-cdn.alpinelinux.org/alpine/edge/testing --no-cache \
    curl gcc musl-dev nim tcc

RUN curl -L https://github.com/nim-lang/Nim/archive/v1.0.2.tar.gz | tar xzf - \
 && curl -L https://github.com/nim-lang/csources/archive/v0.20.0.tar.gz | tar xzf -

RUN mv csources-0.20.0 Nim-1.0.2/csources \
 && cd Nim-*/csources                     \
 && ./build.sh                            \
 && cd ..                                 \
 && ./bin/nim compile -d:release koch     \
 && ./koch boot -d:release                \
 && ./koch install /usr/bin

# TODO Use scratch instead.
RUN apk del --no-cache curl gcc \
 && rm -r                       \
    /Nim-*                      \
    /home                       \
    /media                      \
    /mnt                        \
    /opt                        \
    /root                       \
    /run                        \
    /sbin                       \
    /srv                        \
    /usr/local                  \
    /usr/sbin                   \
    /usr/share                  \
    /var

ENTRYPOINT ["/usr/bin/nim", "--cc:tcc", "--verbosity:0", "-r", "c", "-"]
