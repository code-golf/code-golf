FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base curl

ENV VER=3.9.6

RUN curl -#L https://sourceforge.net/projects/regina-rexx/files/regina-rexx/$VER/regina-rexx-$VER.tar.gz \
  | tar xz --strip-components 1

RUN ./configure && make && make install

COPY rexx.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/rexx -s -static rexx.c

FROM codegolf/lang-base

COPY --from=0 /lib/ld-musl-x86_64.so.1 /lib/
COPY --from=0 /usr/bin/rexx            \
              /usr/local/bin/regina    /usr/bin/
COPY --from=0 /libregina.so.3          /usr/lib/

ENTRYPOINT ["rexx"]

CMD ["--version"]
