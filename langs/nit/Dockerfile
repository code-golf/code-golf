FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base ccache curl gc-dev libunwind-dev pkgconfig

RUN curl -L https://github.com/nitlang/nit/archive/refs/tags/v0.8.tar.gz \
  | tar xz

RUN mv nit-* /usr/local/nit \
 && cd /usr/local/nit       \
 && CFLAGS="-Wno-cpp" make

FROM codegolf/lang-base

COPY --from=0 /bin/cat         \
              /bin/chmod       \
              /bin/mkdir       \
              /bin/sh          /bin/
COPY --from=0 /lib             /lib
COPY --from=0 /usr/bin         /usr/bin
COPY --from=0 /usr/include     /usr/include
COPY --from=0 /usr/lib         /usr/lib
COPY --from=0 /usr/local/nit   /usr/local
COPY --from=0 /usr/libexec/gcc /usr/lib/gcc

COPY /nitwrapper /usr/bin/

RUN chmod +x /usr/bin/nitwrapper

ENTRYPOINT ["nitwrapper"]

CMD ["--version"]
