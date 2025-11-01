FROM debian:bullseye-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc make

ENV HB_BUILD_CONTRIBS='no' HB_BUILD_STRIP='all' HB_USER_LDFLAGS='-static' VER=3.0.0

RUN curl -#L https://downloads.sourceforge.net/harbour-project/harbour-$VER.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc` install

COPY harbour.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/harbour -s -static harbour.c

FROM codegolf/lang-base

COPY --from=0 /usr/bin/harbour     /usr/bin/
COPY --from=0 /usr/local/bin/hbmk2 /usr/local/bin/

ENTRYPOINT ["harbour"]

CMD ["--version"]
