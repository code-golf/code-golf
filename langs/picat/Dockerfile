FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl g++ make

ENV VER=39

RUN curl -# https://picat-lang.org/download/picat${VER}_src.tar.gz | tar xz

WORKDIR /Picat/emu

RUN make -fMakefile.linux64 -j`nproc` \
 && strip picat

COPY picat.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/picat -s /picat.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6          \
              /lib/x86_64-linux-gnu/libgcc_s.so.1      \
              /lib/x86_64-linux-gnu/libm.so.6          \
              /lib/x86_64-linux-gnu/libpthread.so.0    /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2              /lib64/
COPY --from=0 /usr/bin/picat                           /usr/bin/
COPY --from=0 /usr/lib/x86_64-linux-gnu/libstdc++.so.6 /usr/lib/
COPY --from=0 /Picat/emu/picat                         /usr/local/bin/

ENTRYPOINT ["picat"]

CMD ["--version"]
