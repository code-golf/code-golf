FROM debian:bullseye-slim AS builder

RUN apt-get update && apt-get install -y curl gcc

RUN curl https://www.jsoftware.com/download/j9.4/install/j9.4.2_linux64.tar.gz \
  | tar xz

# FIXME Don't run updatejs.sh as it tries to update to AVX512 because my build
# laptop supports it. But prod servers do not. Is there an override flag?
# RUN cd j9.4 && ./updateje.sh

COPY j.c /

RUN gcc -s -o j j.c

FROM codegolf/lang-base

COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgcc_s.so.1   \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /j                                    /usr/bin/
COPY --from=0 /j9.4/addons/dev                      /opt/j/addons/dev
COPY --from=0 /j9.4/bin/jconsole                    \
              /j9.4/bin/libgmp.so                   \
              /j9.4/bin/libj.so                     \
              /j9.4/bin/profile.ijs                 /opt/j/bin/
COPY --from=0 /j9.4/system                          /opt/j/system

ENTRYPOINT ["j"]

CMD ["-v"]
