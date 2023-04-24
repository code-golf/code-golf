FROM debian:bullseye-slim as builder

RUN apt-get update \
 && apt-get install -y curl gcc libncurses-dev libtinfo5 xz-utils

RUN curl -L http://downloads.sourceforge.net/fbc/FreeBASIC-1.09.0-linux-x86_64.tar.xz \
  | tar xJ

RUN cd FreeBASIC-1.09.0-linux-x86_64 && ./install.sh -i /usr

COPY basic.c .

RUN gcc -Wall -Werror -Wextra -o basic -s basic.c

FROM codegolf/lang-base

COPY --from=0 /bin/sh      /bin/
COPY --from=0 /lib         /lib
COPY --from=0 /lib64       /lib64
COPY --from=0 /basic       \
              /usr/bin/as  \
              /usr/bin/fbc \
              /usr/bin/gcc \
              /usr/bin/ld  /usr/bin/
COPY --from=0 /usr/include /usr/include
COPY --from=0 /usr/lib     /usr/lib
COPY --from=0 /usr/libexec /usr/libexec

ENTRYPOINT ["basic"]

CMD ["--version"]
