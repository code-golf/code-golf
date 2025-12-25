FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl g++ make

ENV LDFLAGS='-static' VER=0.30.3

RUN curl -#L https://github.com/IppClub/YueScript/archive/refs/tags/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` install \
 && strip /usr/local/bin/yue

COPY yuescript.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/yuescript -s -static yuescript.c

FROM codegolf/lang-lua

COPY --from=0 /usr/bin/yuescript /usr/bin/
COPY --from=0 /usr/local/bin/yue /usr/local/bin/

ENTRYPOINT ["yuescript"]

CMD ["--version"]
