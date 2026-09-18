FROM debian:trixie-slim AS builder

RUN apt-get update \
 && apt-get install --yes curl gcc unzip binutils

ENV VER=4.33.0

WORKDIR /

RUN curl -sSL -o /lean.zip https://github.com/leanprover/lean4/releases/download/v$VER/lean-$VER-linux.zip \
 && unzip -q lean.zip \
 && mv lean-$VER-linux lean

COPY leanwrapper.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/lean -s -static leanwrapper.c

RUN strip $(find /lean -type f -name "*.a")
RUN strip $(find /lean -type f -name "*.so*" ! -name "*libc++.so" ! -name "*libgcc_s.so")

FROM codegolf/lang-base

# copy everything verbatim for now
COPY --from=0 /lean         /usr/local/lean/
COPY --from=0 /usr/bin/lean /usr/bin/

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                   \
              /lib/x86_64-linux-gnu/libdl.so.2                  \
              /lib/x86_64-linux-gnu/libm.so.6                   \
              /lib/x86_64-linux-gnu/libpthread.so.0             \
              /lib/x86_64-linux-gnu/librt.so.1                  /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                       /lib64/

ENTRYPOINT ["/usr/bin/lean"]

CMD ["--version"]
