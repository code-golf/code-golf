FROM alpine:edge

RUN apk add -X http://dl-cdn.alpinelinux.org/alpine/edge/testing --no-cache \
    curl gcc musl-dev nim tcc

RUN curl -L https://github.com/nim-lang/Nim/archive/16aa10d.tar.gz | tar xzf - \
 && curl -L https://github.com/nim-lang/csources/archive/v0.19.0.tar.gz | tar xzf -

RUN mv csources-0.19.0 Nim-16aa10dfe101da99c402657a420dbf2785ca4a2a/csources \
 && cd Nim-*/csources                      \
 && ./build.sh                             \
 && cd ..                                  \
 && ./bin/nim compile -d:release koch      \
 && ./koch boot -d:release                 \
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
