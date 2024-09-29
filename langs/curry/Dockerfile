FROM swipl:9.2.7 AS builder

RUN apt-get update \
 && apt-get install --yes curl make

ENV LC_ALL="C.UTF-8"

RUN curl https://mirrors.cat.pdx.edu/pakcs/download/pakcs-3.6.0-amd64-Linux.tar.gz | tar xz

RUN mv pakcs* edu \
 && cd edu        \
 && make -j`nproc`

FROM codegolf/lang-base

COPY --from=0 /bin/dirname /bin/readlink /bin/rm /bin/sh    /bin/
COPY --from=0 /edu                                          /edu
COPY --from=0 /lib                                          /lib
COPY --from=0 /lib64/ld-linux-x86-64.so.2                   /lib64/
COPY --from=0 /usr/lib/swipl/bin/x86_64-linux/swipl         /usr/lib/swipl/bin/x86_64-linux/
COPY --from=0 /usr/lib/swipl/lib/x86_64-linux/libswipl.so.9 /usr/lib/swipl/lib/x86_64-linux/

COPY /currywrapper /usr/bin/

ENTRYPOINT ["currywrapper"]

CMD ["--version"]
