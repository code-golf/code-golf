FROM debian:bookworm-slim AS builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
 && apt-get install -y build-essential curl

RUN curl https://www.lua.org/ftp/lua-5.4.7.tar.gz | tar xz --strip-components=1

RUN make clean \
 && make CFLAGS='-fPIC' linux \
 && make install

RUN curl -L https://github.com/pigpigyyy/Yuescript/archive/refs/tags/v0.24.1.tar.gz | tar xz --strip-components=1

RUN make shared \
 && make install

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/ld-linux-x86-64.so.2 \
              /lib/x86_64-linux-gnu/libc.so.6            \
              /lib/x86_64-linux-gnu/libgcc_s.so.1        \
              /lib/x86_64-linux-gnu/libm.so.6            \
              /lib/x86_64-linux-gnu/libstdc++.so.6       /lib/x86_64-linux-gnu/
COPY --from=0 /lib64                                     /lib64
COPY --from=0 /usr/bin/release/yue                       /usr/bin/

ENTRYPOINT ["yue"]

CMD ["-v"]
