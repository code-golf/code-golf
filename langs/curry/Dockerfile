FROM swipl:10.0.0 AS builder

RUN apt-get update                   \
 && DEBIAN_FRONTEND='noninteractive' \
    apt-get install --yes curl gcc gettext make

ENV LC_ALL='C.UTF-8' VER=3.8.0

WORKDIR /usr/local

RUN curl -#L https://www.curry-lang.org/pakcs/download/pakcs-$VER-amd64-Linux.tar.gz \
  | tar xz --strip-components 1

RUN make -j`nproc`

COPY curry.c /

RUN gcc -Wall -Werror -Wextra -o /usr/bin/curry -s -static /curry.c

FROM codegolf/lang-prolog

COPY --from=0 /bin/cat     /bin/chmod    /bin/expr   \
              /bin/grep    /bin/head     /bin/mv     \
              /bin/rm      /bin/sed      /bin/tty    /bin/
COPY --from=0 /bin/dash                              /bin/sh
COPY --from=0 /lib/x86_64-linux-gnu/libacl.so.1      \
              /lib/x86_64-linux-gnu/libattr.so.1     \
              /lib/x86_64-linux-gnu/libdl.so.2       \
              /lib/x86_64-linux-gnu/libpthread.so.0  \
              /lib/x86_64-linux-gnu/libreadline.so.8 \
              /lib/x86_64-linux-gnu/librt.so.1       \
              /lib/x86_64-linux-gnu/libselinux.so.1  \
              /lib/x86_64-linux-gnu/libutil.so.1     /lib/
COPY --from=0 /usr/bin/basename                      \
              /usr/bin/curry                         \
              /usr/bin/dirname                       \
              /usr/bin/readlink                      /usr/bin/
COPY --from=0 /usr/local/pakcsrc.default             /usr/local/
COPY --from=0 /usr/local/bin/curry                   \
              /usr/local/bin/pakcs-fcypp             \
              /usr/local/bin/pakcs-frontend          /usr/local/bin/
COPY --from=0 /usr/local/lib                         /usr/local/lib
COPY --from=0 /usr/local/scripts/makesavedstate      /usr/local/scripts/
COPY --from=0 /usr/local/src/pakcs                   /usr/local/src/

ENTRYPOINT ["curry"]

CMD ["--version"]
