FROM alpine:3.21 AS builder

RUN apk add --no-cache build-base curl

ENV VER=2828c5b

RUN curl -#L https://github.com/themanyone/crap/tarball/$VER \
  | tar xz --strip-components 1 && make && make install

COPY crapwrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/crapwrapper -s crapwrapper.c

FROM codegolf/lang-c

COPY --from=0 /usr/bin/crapwrapper       /usr/bin/
COPY --from=0 /usr/local/bin/crap        /usr/local/bin/
COPY --from=0 /usr/local/include/print.h \
              /usr/local/include/reg.h   \
              /usr/local/include/sjoin.h /usr/local/include/

FROM codegolf/lang-base

COPY --from=1 . /

ENTRYPOINT ["crap"]

CMD ["--dry-run"]
