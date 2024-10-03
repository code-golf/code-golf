FROM debian:bookworm-slim AS builder

ENV VER=4.3.6

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y curl

WORKDIR /haxe

RUN curl -L https://github.com/HaxeFoundation/haxe/releases/download/$VER/haxe-$VER-linux64.tar.gz \
  | tar xz --strip-components 1

FROM codegolf/lang-base

COPY --from=0 /bin/cat /bin/sh              /bin/
COPY --from=0 /lib/x86_64-*/libc.so.6       \
              /lib/x86_64-*/libdl.so.2      \
              /lib/x86_64-*/libm.so.6       \
              /lib/x86_64-*/libpthread.so.0 /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2   /lib64/
COPY --from=0 /haxe/haxe /haxe/haxelib      /usr/bin/
COPY --from=0 /haxe/std                     /usr/lib/std

COPY /haxewrapper /usr/bin/

ENTRYPOINT ["haxewrapper"]

CMD ["--version"]
