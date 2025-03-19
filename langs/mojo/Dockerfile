FROM debian:bookworm-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc

ENV VER=25.1.0

RUN curl -#L https://github.com/modular/mojo/archive/refs/tags/max/v$VER.tar.gz \
  | tar xz --strip-components 1

RUN curl -#L https://magic.modular.com | bash

RUN /root/.modular/bin/magic install \
 && chmod -R 777 /.magic/envs/default

# TOTRY Disable creation on the command line.
RUN touch /.magic/.gitignore

COPY mojo.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/mojo -s -static mojo.c

FROM codegolf/lang-base

# TOFIX There are CRAZY large (and hopefully unneeded) assets in the /.magic/envs/default directory tree.

COPY --from=0 /magic.lock /pixi.toml               /
COPY --from=0 /.magic/.gitignore                   /.magic/
COPY --from=0 /.magic/envs/default/bin             /.magic/envs/default/bin
COPY --from=0 /.magic/envs/default/conda-meta      /.magic/envs/default/conda-meta
COPY --from=0 /.magic/envs/default/etc             /.magic/envs/default/etc
COPY --from=0 /.magic/envs/default/lib             /.magic/envs/default/lib
COPY --from=0 /.magic/envs/default/share           /.magic/envs/default/share
COPY --from=0 /bin                                 /bin
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6      \
              /lib/x86_64-linux-gnu/libdl.so.2     \
              /lib/x86_64-linux-gnu/libgcc_s.so.1  \
              /lib/x86_64-linux-gnu/libm.so.6      \
              /lib/x86_64-linux-gnu/librt.so.1     \
              /lib/x86_64-linux-gnu/libstdc++.so.6 \
              /lib/x86_64-linux-gnu/libtinfo.so.6  \
              /lib/x86_64-linux-gnu/libz.so.1      /lib/x86_64-linux-gnu/
COPY --from=0 /lib64/ld-linux-x86-64.so.2          /lib64/
COPY --from=0 /root/.modular/bin/magic             \
              /usr/bin/env                         \
              /usr/bin/mojo                        /usr/bin/

ENTRYPOINT ["mojo"]

CMD ["--version"]
