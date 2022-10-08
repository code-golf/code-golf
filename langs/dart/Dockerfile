FROM dart:2.18.2-sdk as builder

RUN apt-get update && apt-get install -y gcc

COPY dart-wrapper.c /

RUN cd / && gcc -s -o dart-wrapper dart-wrapper.c

FROM codegolf/lang-base

COPY --from=0 /dart-wrapper /usr/bin/
COPY --from=0 /bin          /bin
COPY --from=0 /lib          /lib
COPY --from=0 /lib64        /lib64
COPY --from=0 /usr/bin      /usr/bin
COPY --from=0 /usr/include  /usr/include
COPY --from=0 /usr/lib      /usr/lib
COPY --from=0 /usr/libexec  /usr/libexec

ENTRYPOINT ["dart-wrapper"]

CMD ["--version"]
