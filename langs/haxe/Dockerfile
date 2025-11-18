FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc

ENV VER=4.3.7

RUN curl -#L https://github.com/HaxeFoundation/haxe/releases/download/$VER/haxe-$VER-linux64.tar.gz \
  | tar xz --directory /usr/local --strip-components 1

COPY haxe.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/haxe -s haxe.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/haxe                         /usr/bin/
COPY --from=0 /usr/local/haxe                       \
              /usr/local/haxelib                    /usr/local/bin/
COPY --from=0 /usr/local/std                        /usr/local/lib/std

ENTRYPOINT ["haxe"]

CMD ["--version"]
