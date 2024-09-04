FROM debian:bookworm-slim AS builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
 && apt-get install -y curl

RUN curl -sSL https://github.com/koka-lang/koka/releases/download/v3.1.2/install.sh | sh

FROM codegolf/lang-base

COPY --from=0 / /

COPY --from=0 /bin                                  /bin
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6       \
              /lib/x86_64-linux-gnu/libdl.so.2      \
              /lib/x86_64-linux-gnu/libgmp.so.10    \
              /lib/x86_64-linux-gnu/libm.so.6       \
              /lib/x86_64-linux-gnu/libpthread.so.0 \
              /lib/x86_64-linux-gnu/librt.so.1      \
              /lib/x86_64-linux-gnu/libutil.so.1    /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2           /lib64/
COPY --from=0 /usr/bin/gcc /usr/bin/ld              /usr/bin/
COPY --from=0 /usr/local/bin                        /usr/local/bin
COPY --from=0 /usr/local/lib                        /usr/local/lib
COPY --from=0 /usr/local/share                      /usr/local/share

COPY /kokawrapper /usr/bin/
RUN chmod +x /usr/bin/kokawrapper

ENTRYPOINT ["kokawrapper"]

CMD ["--version"]
