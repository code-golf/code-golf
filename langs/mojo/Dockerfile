FROM debian:trixie-slim AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes gcc wget

ENV VER=25.6.0

RUN wget -qO- https://pixi.sh/install.sh \
  | PIXI_ARCH='x86_64' PIXI_HOME='/usr' sh

RUN pixi init -c https://conda.modular.com/max -c conda-forge

RUN pixi add modular==$VER

COPY mojo.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/mojo -s mojo.c

FROM codegolf/lang-base

COPY --from=0 /pixi.lock /pixi.toml                                       /
COPY --from=0 /.pixi/.condapackageignore /.pixi/.gitignore                /.pixi/
COPY --from=0 /.pixi/envs/default/bin/modular-crashpad-handler            \
              /.pixi/envs/default/bin/mojo                                /.pixi/envs/default/bin/
COPY --chown=nobody                                                       \
     --from=0 /.pixi/envs/default/conda-meta/*.json                       \
              /.pixi/envs/default/conda-meta/history                      \
              /.pixi/envs/default/conda-meta/pixi_env_prefix              /.pixi/envs/default/conda-meta/
COPY --from=0 /.pixi/envs/default/etc/conda/activate.d/10-activate-max.sh /.pixi/envs/default/etc/conda/activate.d/
COPY --from=0 /.pixi/envs/default/lib/libAsyncRTMojoBindings.so           \
              /.pixi/envs/default/lib/libAsyncRTRuntimeGlobals.so         \
              /.pixi/envs/default/lib/libKGENCompilerRTShared.so          \
              /.pixi/envs/default/lib/libMSupportGlobals.so               \
              /.pixi/envs/default/lib/libNVPTX.so                         /.pixi/envs/default/lib/
COPY --from=0 /.pixi/envs/default/lib/mojo/stdlib.mojopkg                 /.pixi/envs/default/lib/mojo/
COPY --chown=nobody                                                       \
     --from=0 /.pixi/envs/default/share/max/modular.cfg                   /.pixi/envs/default/share/max/
COPY --from=0 /bin/bash /bin/ldd                                          /bin/
COPY --from=0 /lib/x86_64-linux-gnu/libc.so.6                             \
              /lib/x86_64-linux-gnu/libdl.so.2                            \
              /lib/x86_64-linux-gnu/libgcc_s.so.1                         \
              /lib/x86_64-linux-gnu/libm.so.6                             \
              /lib/x86_64-linux-gnu/librt.so.1                            \
              /lib/x86_64-linux-gnu/libstdc++.so.6                        \
              /lib/x86_64-linux-gnu/libtinfo.so.6                         /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                                 /lib64/
COPY --from=0 /usr/bin/env /usr/bin/mojo /usr/bin/pixi                    /usr/bin/

ENTRYPOINT ["mojo"]

CMD ["--version"]
