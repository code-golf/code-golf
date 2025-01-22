FROM debian:bullseye-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes cmake curl

ENV VER=b7bf640

RUN curl -#L https://github.com/VorpalBlade/cfunge/tarball/$VER \
  | tar xz --strip-components 1

WORKDIR /build

RUN cmake .. && make install \
 && strip /usr/local/bin/cfunge

COPY befunge.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/befunge -s /befunge.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libncurses.so.6 \
              /lib/x86_64-linux-gnu/libtinfo.so.6   /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/befunge                      \
              /usr/local/bin/cfunge                 /usr/bin/

ENTRYPOINT ["befunge"]

CMD ["--version"]
