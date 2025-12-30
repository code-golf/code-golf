FROM alpine:3.23 AS builder

RUN apk add --no-cache curl gcc make musl-dev

ENV VER=2828c5b

RUN curl -#L https://github.com/themanyone/crap/tarball/$VER \
  | tar xz --strip-components 1

RUN make \
 && make install

COPY crap.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/crap -s crap.c

FROM codegolf/lang-c

COPY --from=0 /usr/bin/crap              /usr/bin/
COPY --from=0 /usr/local/bin/crap        /usr/local/bin/
COPY --from=0 /usr/local/include/print.h \
              /usr/local/include/reg.h   \
              /usr/local/include/sjoin.h /usr/local/include/

ENTRYPOINT ["crap"]
