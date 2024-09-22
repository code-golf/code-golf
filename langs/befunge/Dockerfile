FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base cmake curl ncurses-dev

RUN curl -L https://github.com/VorpalBlade/cfunge/archive/refs/tags/1,001.tar.gz \
  | tar xz --strip-components 1

RUN cd tmp   \
 && cmake .. \
 && make

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh          /bin/
COPY --from=0 /lib/ld-musl-x86_64.so.1  /lib/
COPY --from=0 /tmp/cfunge               /usr/bin/
COPY --from=0 /usr/lib/libncursesw.so.6 /usr/lib/

COPY /befunge /usr/bin/

ENTRYPOINT ["befunge"]

CMD ["-V"]
