FROM currylang/pakcs:3.6.0 AS builder

RUN apt-get remove -y rlwrap

RUN cypm update \
 && cypm install runcurry

FROM codegolf/lang-base

COPY --from=0 /bin                                       /bin
COPY --from=0 /lib/x86_64-linux-gnu/libacl.so.1          \
              /lib/x86_64-linux-gnu/libattr.so.1         \
              /lib/x86_64-linux-gnu/libc.so.6            \
              /lib/x86_64-linux-gnu/libdl.so.2           \
              /lib/x86_64-linux-gnu/libgmp.so.10         \
              /lib/x86_64-linux-gnu/libm.so.6            \
              /lib/x86_64-linux-gnu/libpcre.so.3         \
              /lib/x86_64-linux-gnu/libpcre2-8.so.0      \
              /lib/x86_64-linux-gnu/libpthread.so.0      \
              /lib/x86_64-linux-gnu/libreadline.so.8     \
              /lib/x86_64-linux-gnu/librt.so.1           \
              /lib/x86_64-linux-gnu/libselinux.so.1      \
              /lib/x86_64-linux-gnu/libtinfo.so.6        \
              /lib/x86_64-linux-gnu/libutil.so.1         \
              /lib/x86_64-linux-gnu/libz.so.1            /lib/
COPY --from=0 /lib64/ld-linux-x86-64.so.2                /lib64/
COPY --from=0 /pakcs/pakcs                               /pakcs/pakcs
COPY --from=0 /userhome/.cpm/bin/runcurry                \
              /usr/bin/basename                          \
              /usr/bin/dirname                           \
              /usr/bin/readlink                          \
              /usr/bin/tty                               /usr/bin/
COPY --from=0 /usr/lib/libswipl.so.8                     /usr/lib/
COPY --from=0 /usr/lib/swi-prolog/bin/x86_64-linux/swipl /usr/lib/swi-prolog/bin/x86_64-linux/
COPY --from=0 /usr/lib/swi-prolog/lib                    /usr/lib/swi-prolog/lib

COPY /currywrapper /usr/bin/

ENV PATH=/pakcs/pakcs/bin:$PATH

ENTRYPOINT ["currywrapper"]

CMD ["--version"]
