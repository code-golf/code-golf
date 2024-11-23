FROM alpine:3.20 AS builder

RUN wget -cO /usr/bin/iogii https://golfscript.com/iogii/iogii-0.1-alpha \
 && chmod +x /usr/bin/iogii

FROM codegolf/lang-ruby AS interpreter

COPY --from=0 /bin/cat /bin/sh /bin/touch /bin/
COPY --from=0 /usr/bin/env /usr/bin/iogii /usr/bin/

FROM codegolf/lang-base

COPY --from=1 . /

COPY iogiiwrapper /usr/bin/

ENTRYPOINT ["iogiiwrapper"]

CMD ["--version"]
