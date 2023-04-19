FROM debian:bullseye-slim as builder

RUN apt-get update \
 && apt-get install -y curl gcc libncurses-dev libtinfo5 xz-utils

RUN curl -L http://downloads.sourceforge.net/fbc/FreeBASIC-1.09.0-linux-x86_64.tar.xz \
  | tar xJ

RUN cd FreeBASIC-1.09.0-linux-x86_64 && ./install.sh -i /usr

FROM codegolf/lang-base

COPY --from=0 /bin         /bin
COPY --from=0 /lib         /lib
COPY --from=0 /lib64       /lib64
COPY --from=0 /usr/bin/as  \
              /usr/bin/fbc \
              /usr/bin/gcc \
              /usr/bin/ld  /usr/bin/
COPY --from=0 /usr/include /usr/include
COPY --from=0 /usr/lib     /usr/lib
COPY --from=0 /usr/libexec /usr/libexec

COPY basic /usr/bin/

ENTRYPOINT ["basic"]

CMD ["--version"]
